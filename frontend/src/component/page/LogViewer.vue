<template>
  <div class="flex flex-col flex-1 overflow-hidden">
    <!-- Toolbar -->
    <div class="p-2 space-y-2">
      <div class="flex items-center gap-2">
        <button class="px-3 py-1.5 cursor-pointer rounded-md border border-stone-600 bg-stone-700/60 text-stone-100 text-sm hover:bg-stone-700"
                @click="openFile">
          開啟檔案
        </button>

        <select class="px-2 py-1 rounded-md border border-stone-700 bg-stone-900 text-stone-100 text-sm w-44"
                v-model="selectedFromDir" @change="chooseFromDir">
          <option value="">從 log 目錄選取...</option>
          <option v-for="file in files" :key="file.path" :value="file.path">
            {{ file.name }}
          </option>
        </select>

        <span v-if="filename"
              :title="filename"
              class="ml-auto inline-flex items-center px-2 py-0.5 rounded border border-stone-700 bg-stone-800 text-xs truncate max-w-[60%]">
          {{ filename }}
        </span>
      </div>

      <div class="flex flex-wrap items-center gap-2">
        <input v-model="search" placeholder="搜尋..." />
        <input v-model="filter" placeholder="篩選..." />
        <input v-model="highlight" placeholder="標記..." />

        <div class="ml-auto flex items-center gap-2">
          <button class="cursor-pointer px-3 py-1.5 rounded-md border border-stone-600 text-stone-100 text-sm hover:bg-stone-700/50 disabled:opacity-50"
                  @click="loadMore" :disabled="loading || atFileEnd || !filename">
            載入更多
          </button>
          <button class="cursor-pointer px-3 py-1.5 rounded-md border border-stone-600 text-stone-100 text-sm hover:bg-stone-700/50"
                  @click="scrollBottom">
            跳至底部
          </button>
        </div>
      </div>
    </div>

    <!-- 表頭 -->
    <div class="sticky top-0 z-10 bg-stone-950">
      <div class="grid grid-cols-[80px_1fr] items-center font-semibold text-slate-200 text-xs px-2 py-1">
        <div>#</div>
        <div>內容</div>
      </div>
    </div>

    <!-- 可變高度虛擬清單（由頭→尾） -->
    <DynamicScroller
      ref="scrollerRef"
      class="flex-1 overflow-auto"
      :items="displayedRows"
      :min-item-size="rowMinHeight"
      key-field="id"
      @scroll="onVsScroll"
    >
      <template #default="{ item, index, active }">
        <DynamicScrollerItem :item="item" :active="active" :data-index="index">
          <div class="grid grid-cols-[80px_1fr] items-start font-mono text-xs border-b border-dashed border-stone-700/70 px-2 py-1">
            <div class="text-slate-400 select-none">{{ item.line }}</div>
            <div class="text-slate-200 text-left whitespace-pre-wrap break-words" v-html="renderMsg(item.content)"></div>
          </div>
        </DynamicScrollerItem>
      </template>
    </DynamicScroller>

    <!-- Footer -->
    <div class="flex items-center justify-between text-xs text-slate-300 bg-stone-950 px-2 py-1.5">
      <div>已載入行數：{{ rows.length }} / 總行數：{{ totalLines }}</div>
      <div v-if="loading">載入中...</div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import { DynamicScroller, DynamicScrollerItem } from 'vue3-virtual-scroller'
import 'vue3-virtual-scroller/dist/vue3-virtual-scroller.css'

import {
  ListLogFiles,
  ReadFirstLines,   // 由頭開始（依行數）
  LoadNext,         // 往後續讀（依行數）
  SetCurrentFile,
  PickFile,
  CountLines,
} from '@/wailsjs/go/log_viewer/LogViewerService.js'

/* ------------ 常數：每次載入行數 ------------ */
const PAGE_LINES = 800

/* ------------ 狀態 ------------ */
const filename = ref('')
const files = ref([])          // [{ name, path }]
const selectedFromDir = ref('')

const rows = ref([])           // [{ id, line, content, raw }]
let nextId = 1                 // 穩定遞增 id
const maxKeep = 200000
const totalLines = ref(0)      // 整檔總行數

const anchor = ref(0)          // 下一次往後讀的位移（由後端回傳）
const atFileEnd = ref(false)   // 是否已到檔案尾
const loading = ref(false)

const search = ref('')
const filter = ref('')
const highlight = ref('')

const rowMinHeight = 22
const scrollerRef = ref(null)
let loadingLock = false

/* ------------ 顯示資料（過濾） ------------ */
const displayedRows = computed(() => {
  let list = rows.value
  const s = (search.value || '').toLowerCase()
  const f = (filter.value || '').toLowerCase()
  if (s) list = list.filter(r => (r.content || '').toLowerCase().includes(s))
  if (f) list = list.filter(r => (r.content || '').toLowerCase().includes(f))
  return list
})

/* ------------ lifecycle ------------ */
onMounted(async () => {
  await fetchDirFiles()
})

/* ------------ Scroll handlers：靠近底部就載入下一塊 ------------ */
function onVsScroll() {
  const el = scrollerRef.value?.$el
  if (!el) return
  const nearBottom = el.scrollHeight - (el.scrollTop + el.clientHeight) <= 60
  if (nearBottom && !loading.value && !atFileEnd.value && !loadingLock) {
    loadingLock = true
    const wasAtBottom = nearBottom
    loadMore().finally(() => {
      if (wasAtBottom) scrollBottom()
      loadingLock = false
    })
  }
}

/* ------------ 文字處理 / 高亮 ------------ */
function esc(s) {
  return (s || '').replace(/[&<>"']/g, m => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[m]))
}
function renderMsg(msg) {
  if (!highlight.value) return esc(msg)
  const key = esc(highlight.value)
  const re = new RegExp(key, 'gi')
  return esc(msg).replace(re, m => `<span class="bg-yellow-500/40 rounded px-0.5">${m}</span>`)
}

/* ------------ 後端互動 ------------ */
async function openFile() {
  try {
    const path = await PickFile()
    await useFile(path)
  } catch (_) {}
}
async function fetchDirFiles() {
  try { files.value = await ListLogFiles() } catch (e) { console.error(e) }
}
async function chooseFromDir() {
  if (!selectedFromDir.value) return
  await useFile(selectedFromDir.value)
}
async function useFile(path) {
  filename.value = path
  rows.value = []
  nextId = 1
  anchor.value = 0
  atFileEnd.value = false

  try { await SetCurrentFile(path) } catch (_) {}

  try { totalLines.value = await CountLines(path) } catch { totalLines.value = 0 }

  await loadInitial()
  nextTick(scrollBottom)
}

async function loadInitial() {
  if (!filename.value) return
  loading.value = true
  try {
    const res = await ReadFirstLines(filename.value, PAGE_LINES)
    anchor.value = res.nextStart
    atFileEnd.value = res.eof
    replaceRows(res.lines)   // 初次載入：整批替換
  } finally {
    loading.value = false
  }
}

async function loadMore() {
  if (!filename.value || atFileEnd.value || loading.value) return
  loading.value = true
  try {
    const res = await LoadNext(filename.value, anchor.value, PAGE_LINES)
    anchor.value = res.nextStart
    atFileEnd.value = res.eof
    pushRows(res.lines)      // 追加到底部
  } finally {
    loading.value = false
  }
}

/* ------------ rows 操作 ------------ */
function makeRow(line, index) {
  return { id: nextId++, line: index + 1, content: line, raw: line }
}
function renumberRows() {
  for (let i = 0; i < rows.value.length; i++) rows.value[i].line = i + 1
}
function replaceRows(lines) {
  const parsed = (lines || []).filter(Boolean).map((l, idx) => makeRow(l, idx))
  rows.value = parsed
  renumberRows()
  trimRows()
}
function pushRows(lines) {
  const parsed = (lines || []).filter(Boolean).map((l) => ({ id: nextId++, line: 0, content: l, raw: l }))
  rows.value = rows.value.concat(parsed)
  renumberRows()
}
function trimRows() {
  if (rows.value.length > maxKeep) {
    const cut = rows.value.length - maxKeep
    rows.value.splice(0, cut)
    renumberRows()
  }
}
function scrollBottom() {
  const el = scrollerRef.value?.$el
  if (!el) return
  el.scrollTop = el.scrollHeight
}
</script>

<style scoped>
@reference "@/src/css/style.css";

input {
  @apply px-3 py-1.5 rounded-md border border-stone-700 bg-stone-900 text-stone-100 text-sm min-w-[12rem] flex-1;
}
</style>