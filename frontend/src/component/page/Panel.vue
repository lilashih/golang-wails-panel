<template>
  <Layout :title="'應用程式'" :items="projects">
    <template #default="{ item, index }">
      <div class="truncate w-1/3">{{ item.Config.title }}</div>
      <div class="flex justify-end items-center gap-4 w-2/3">
        <div
          class="status"
          :class="item.running ? 'running' : 'stopped'"
        >
          {{ item.running ? '啟動中' : '已停止' }}
        </div>
        <div class="flex gap-2">
          <!-- 啟動／停止二選一 -->
          <button
            v-if="!item.running"
            class="btn btn-start"
            @click="startProject(index)"
            :disabled="loading"
          >
            啟動
          </button>
          <button
            v-else
            class="btn btn-stop"
            @click="stopProject(index)"
            :disabled="loading"
          >
            停止
          </button>
          <!-- 非預期狀態也能手動停止 -->
          <button
            class="btn btn-stop"
            @click="stopProject(index)"
            :disabled="loading"
          >
            強制停止
          </button>
          <!-- 重新安裝 -->
          <button
            class="btn"
            @click="installProject(index)"
            :disabled="loading || item.running"
          >
            重新安裝
          </button>
          <!-- 開啟 -->
          <button
            class="btn"
            @click="openProject(index)"
            :disabled="loading"
          >
            開啟
          </button>
        </div>
      </div>
    </template>
  </Layout>
</template>

<script setup>
import Layout from '@/src/component/layout/LoggerLayout.vue'
import { ref, onMounted } from 'vue'
import { useLoadingStore } from '@/src/store/loading.js'
import {
  ListProjects,
  OpenProject,
  StartProject,
  StopProject,
  InstallProject,
} from '@/wailsjs/go/panel/PanelService.js'


/* ---------- 狀態 ---------- */
const projects = ref([])
const loadingStore = useLoadingStore()

/* ---------- 讀取 ---------- */
async function load() {
  loadingStore.setLoading(true, '載入應用程式中...')
  projects.value = await ListProjects()
  loadingStore.setLoading(false)
}

/* ---------- 操作 ---------- */
async function openProject(index) {
  loadingStore.setLoading(true, '')
  await OpenProject(index)

  setTimeout(() => {
    loadingStore.setLoading(false)
  }, 2000)
}

async function startProject(index) {
  loadingStore.setLoading(true, '啟動中，請稍候...')
  projects.value[index].running = await StartProject(index)
  loadingStore.setLoading(false)
}

async function stopProject(index) {
  loadingStore.setLoading(true, '停止中，請稍候...')
  projects.value[index].running = await StopProject(index)
  loadingStore.setLoading(false)
}

async function installProject(index) {
  loadingStore.setLoading(true, '安裝中，請稍候...')
  await InstallProject(index)
  loadingStore.setLoading(false)
}

onMounted(load)
</script>

<style scoped>
@reference "@/src/css/style.css";

.status {
  @apply text-xs px-1 py-0.5 rounded border
}
.running {
  @apply text-teal-400 border-green-700
}
.stopped {
  @apply text-gray-400 border-gray-600
}
</style>