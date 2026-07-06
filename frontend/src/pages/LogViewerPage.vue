<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from "vue";
import { useLoadingStore } from "@/stores/loading";
import { Events } from "@wailsio/runtime";
import { DynamicScroller, DynamicScrollerItem } from "vue3-virtual-scroller";
import { ChevronsDown, FileText, Filter, Highlighter, RefreshCw, Search, ToggleLeft, ToggleRight } from "lucide-vue-next";
import {
  CountLines,
  ListLogFiles,
  LoadNext,
  OpenLogFileLocation,
  ReadFirstLines,
  SetCurrentFile,
  StartTail,
  StopTail,
} from "@bindings/gbase/src/service/log_viewer/service.js";
import AppButton from "@/components/ui/AppButton.vue";
import EmptyState from "@/components/ui/EmptyState.vue";
import IconButton from "@/components/ui/IconButton.vue";

interface LogFileItem {
  name: string;
  path: string;
}

interface LogRow {
  id: number;
  text: string;
}

interface TailEvent {
  path: string;
  lines: string[];
  timestamp: string;
}

const PAGE_LINES = 800;
const MAX_LOG_ROWS = 10000;

const files = ref<LogFileItem[]>([]);
const currentPath = ref("");
const rows = ref<LogRow[]>([]);
const nextStart = ref(0);
const eof = ref(true);
const totalLines = ref(0);
const search = ref("");
const filter = ref("");
const highlight = ref("");
const followTail = ref(false);
const loading = ref(false);
const error = ref("");
const scrollerRef = ref<InstanceType<typeof DynamicScroller> | null>(null);
const loadingStore = useLoadingStore();
let offTail: (() => void) | null = null;
let rowId = 0;

const filteredRows = computed(() => {
  const searchKeyword = search.value.trim().toLowerCase();
  const filterKeyword = filter.value.trim().toLowerCase();
  if (!searchKeyword && !filterKeyword) {
    return rows.value;
  }

  return rows.value.filter((row) => {
    const text = row.text.toLowerCase();
    return (!searchKeyword || text.includes(searchKeyword)) && (!filterKeyword || text.includes(filterKeyword));
  });
});

function toErrorMessage(err: unknown) {
  return err instanceof Error ? err.message : String(err);
}

function isLogFileItem(value: unknown): value is LogFileItem {
  if (!value || typeof value !== "object") {
    return false;
  }

  const file = value as Partial<LogFileItem>;
  return typeof file.name === "string" && typeof file.path === "string";
}

function normalizeLogFiles(value: unknown): LogFileItem[] {
  if (!Array.isArray(value)) {
    return [];
  }

  return value.filter(isLogFileItem);
}

function normalizeLines(value: unknown): string[] {
  if (!Array.isArray(value)) {
    return [];
  }

  return value.filter((line): line is string => typeof line === "string");
}

function isTailEvent(value: unknown): value is TailEvent {
  if (!value || typeof value !== "object") {
    return false;
  }

  const data = value as Partial<TailEvent>;
  return typeof data.path === "string" && Array.isArray(data.lines) && data.lines.every((line) => typeof line === "string");
}

function appendLines(lines: string[]) {
  if (lines.length === 0) {
    return;
  }

  rows.value.push(...lines.map((line) => ({ id: rowId++, text: line })));
  const overflow = rows.value.length - MAX_LOG_ROWS;
  if (overflow > 0) {
    rows.value.splice(0, overflow);
  }
}

function scrollToBottom() {
  nextTick(() => {
    const scroller = scrollerRef.value as unknown as { scrollToItem?: (index: number) => void } | null;
    if (!scroller || filteredRows.value.length === 0) {
      return;
    }
    scroller.scrollToItem?.(filteredRows.value.length - 1);
  });
}

watch(followTail, (enabled) => {
  if (enabled) {
    scrollToBottom();
  }
});

async function loadFiles() {
  loading.value = true;
  error.value = "";
  try {
    await loadingStore.withLoading("讀取中", async () => {
      files.value = normalizeLogFiles(await ListLogFiles());
      if (!currentPath.value && files.value.length > 0) {
        await selectFile(files.value[files.value.length - 1].path);
      }
    });
  } catch (err) {
    error.value = toErrorMessage(err);
  } finally {
    loading.value = false;
  }
}

async function selectFile(path: string) {
  loading.value = true;
  error.value = "";
  try {
    if (followTail.value) {
      await StopTail();
      followTail.value = false;
    }
    currentPath.value = await SetCurrentFile(path);
    const result = await ReadFirstLines(currentPath.value, PAGE_LINES);
    rows.value = [];
    rowId = 0;
    appendLines(normalizeLines(result?.lines));
    nextStart.value = result?.nextStart || 0;
    eof.value = Boolean(result?.eof);
    totalLines.value = await CountLines(currentPath.value);
    scrollToBottom();
  } catch (err) {
    error.value = toErrorMessage(err);
  } finally {
    loading.value = false;
  }
}

async function loadMore() {
  if (!currentPath.value || eof.value) {
    return;
  }

  loading.value = true;
  error.value = "";
  try {
    await loadingStore.withLoading("載入中", async () => {
      const result = await LoadNext(currentPath.value, nextStart.value, PAGE_LINES);
      appendLines(normalizeLines(result?.lines));
      nextStart.value = result?.nextStart || nextStart.value;
      eof.value = Boolean(result?.eof);
      if (followTail.value) {
        scrollToBottom();
      }
    });
  } catch (err) {
    error.value = toErrorMessage(err);
  } finally {
    loading.value = false;
  }
}

async function toggleTail() {
  if (!currentPath.value) {
    return;
  }

  error.value = "";
  try {
    if (followTail.value) {
      await StopTail();
      followTail.value = false;
      return;
    }

    await StartTail(currentPath.value);
    followTail.value = true;
    scrollToBottom();
  } catch (err) {
    error.value = toErrorMessage(err);
  }
}

function onTail(event: { data: unknown }) {
  if (!followTail.value || !isTailEvent(event.data) || event.data.path !== currentPath.value) {
    return;
  }

  appendLines(event.data.lines);
  totalLines.value += event.data.lines.length;
  scrollToBottom();
}

function escapeHtml(value: string) {
  const entities: Record<string, string> = {
    "&": "&amp;",
    "<": "&lt;",
    ">": "&gt;",
    '"': "&quot;",
    "'": "&#39;",
  };
  return value.replace(/[&<>"']/g, (char) => entities[char] || char);
}

function escapeRegExp(value: string) {
  return value.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
}

function renderLogText(text: string) {
  const escaped = escapeHtml(text || "");
  const keyword = highlight.value.trim();
  if (!keyword) {
    return escaped;
  }

  const pattern = new RegExp(escapeRegExp(escapeHtml(keyword)), "gi");
  return escaped.replace(pattern, (match) => `<mark class="rounded bg-app-warning-soft px-0.5 text-app-text">${match}</mark>`);
}

async function openCurrentFileLocation() {
  if (!currentPath.value) {
    return;
  }

  error.value = "";
  try {
    await loadingStore.withLoading("開啟檔案中", () => OpenLogFileLocation(currentPath.value));
  } catch (err) {
    error.value = toErrorMessage(err);
  }
}

onMounted(() => {
  loadFiles().catch(console.error);
  offTail = Events.On("logViewer:tail", onTail);
});

onBeforeUnmount(() => {
  offTail?.();
  offTail = null;

  if (followTail.value) {
    StopTail().catch(console.error);
  }
});
</script>

<template>
  <section class="flex h-full min-h-0 flex-col overflow-hidden">
    <div class="flex-none border-b border-app-border bg-app-surface px-4 py-3">
      <div class="mb-3 flex items-center justify-between gap-3">
        <p class="text-sm text-app-text-muted">
          已載入 {{ rows.length }} 行<span v-if="totalLines"> / 共 {{ totalLines }} 行</span>
        </p>
        <div class="flex items-center gap-2 text-sm">
          <AppButton size="sm" @click="loadFiles">
            <RefreshCw class="h-4 w-4" />
            重新整理
          </AppButton>
          <AppButton size="sm" :variant="followTail ? 'primary' : 'secondary'" :disabled="!currentPath" @click="toggleTail">
            <ToggleRight v-if="followTail" class="h-4 w-4" />
            <ToggleLeft v-else class="h-4 w-4" />
            tail follow
          </AppButton>
        </div>
      </div>

      <div class="grid grid-cols-[15rem_repeat(3,minmax(8rem,1fr))_auto] gap-3 text-sm">
        <label class="relative">
          <FileText class="pointer-events-none absolute left-3 top-2.5 h-4 w-4 text-app-text-subtle" />
          <select
            v-model="currentPath"
            class="h-9 w-full rounded-md border border-app-border bg-app-surface pl-9 pr-8 text-sm text-app-text focus:border-app-primary focus:outline-none focus:ring-2 focus:ring-app-primary/20"
            @change="selectFile(currentPath)"
          >
            <option value="">選擇 log 檔案</option>
            <option v-for="file in files" :key="file.path" :value="file.path">{{ file.name }}</option>
          </select>
        </label>

        <label class="relative">
          <Search class="pointer-events-none absolute left-3 top-2.5 h-4 w-4 text-app-text-subtle" />
          <input
            v-model="search"
            type="search"
            class="h-9 w-full rounded-md border border-app-border bg-app-surface pl-9 pr-3 text-sm text-app-text focus:border-app-primary focus:outline-none focus:ring-2 focus:ring-app-primary/20"
            placeholder="搜尋"
          />
        </label>

        <label class="relative">
          <Filter class="pointer-events-none absolute left-3 top-2.5 h-4 w-4 text-app-text-subtle" />
          <input
            v-model="filter"
            type="search"
            class="h-9 w-full rounded-md border border-app-border bg-app-surface pl-9 pr-3 text-sm text-app-text focus:border-app-primary focus:outline-none focus:ring-2 focus:ring-app-primary/20"
            placeholder="篩選"
          />
        </label>

        <label class="relative">
          <Highlighter class="pointer-events-none absolute left-3 top-2.5 h-4 w-4 text-app-text-subtle" />
          <input
            v-model="highlight"
            type="search"
            class="h-9 w-full rounded-md border border-app-border bg-app-surface pl-9 pr-3 text-sm text-app-text focus:border-app-primary focus:outline-none focus:ring-2 focus:ring-app-primary/20"
            placeholder="標記"
          />
        </label>

        <IconButton label="跳到底部" :disabled="rows.length === 0" @click="scrollToBottom">
          <ChevronsDown class="h-4 w-4" />
        </IconButton>
      </div>

      <div v-if="error" class="mt-3 rounded-md border border-app-danger/30 bg-app-danger-soft px-3 py-2 text-sm text-app-danger">
        {{ error }}
      </div>
    </div>

    <div class="min-h-0 flex-1 overflow-hidden p-4">
      <EmptyState
        v-if="!currentPath && files.length === 0 && !loading"
        title="目前沒有 log 檔"
        message="後端 logger 寫入 storage/log 後，這裡會顯示可選擇的檔案。"
      />

      <div v-else class="flex h-full min-h-0 flex-col overflow-hidden rounded-md border border-app-border bg-app-surface">
        <div class="flex-none border-b border-app-border bg-app-surface-muted px-3 py-2 text-xs text-app-text-muted">
          <div class="flex items-center justify-between gap-3">
            <button
              type="button"
              class="min-w-0 truncate rounded px-1 text-left font-mono hover:bg-app-bg hover:text-app-text disabled:cursor-default disabled:hover:bg-transparent disabled:hover:text-app-text-muted"
              :class="currentPath ? 'cursor-pointer' : 'text-app-text-muted'"
              :disabled="!currentPath"
              :title="currentPath ? '開啟檔案位置：' + currentPath : '未選擇檔案'"
              @click="openCurrentFileLocation"
            >
              {{ currentPath || "未選擇檔案" }}
            </button>
            <span>{{ eof ? "EOF" : "尚可載入更多" }}</span>
          </div>
        </div>

        <div class="min-h-0 flex-1 overflow-hidden bg-app-log-bg font-mono text-xs leading-5 text-app-text-inverse">
          <DynamicScroller
            ref="scrollerRef"
            class="h-full"
            :items="filteredRows"
            :min-item-size="28"
            key-field="id"
            v-slot="{ item, index, active }"
          >
            <DynamicScrollerItem
              :item="item"
              :active="active"
              :size-dependencies="[item.text, highlight]"
              :data-index="index"
            >
              <div class="grid min-h-7 grid-cols-[4.5rem_minmax(0,1fr)] border-b border-app-log-divider/5">
                <span class="select-none border-r border-app-log-divider/10 px-3 py-1 text-right text-app-text-muted mono-tabular">
                  {{ index + 1 }}
                </span>
                <span class="break-words whitespace-pre-wrap px-3 py-1" v-html="renderLogText(item.text)"></span>
              </div>
            </DynamicScrollerItem>
          </DynamicScroller>
        </div>

        <footer class="flex h-10 flex-none items-center justify-between border-t border-app-border bg-app-surface-muted px-3 text-xs text-app-text-muted">
          <span>顯示 {{ filteredRows.length }} / 已載入 {{ rows.length }}</span>
          <AppButton size="sm" :disabled="eof || loading || !currentPath" @click="loadMore">
            載入下一批
          </AppButton>
        </footer>
      </div>
    </div>
  </section>
</template>
