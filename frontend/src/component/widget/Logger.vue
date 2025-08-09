<template>
  <div class="flex flex-col h-full bg-neutral-700 whitespace-pre-wrap break-all rounded-lg">
    <div class="flex items-center justify-between px-4 py-2 bg-zinc-950 rounded-t-lg">
      <h3 class="font-semibold">{{ title }}</h3>
      <button
        class="cursor-pointer text-sm text-gray-500 hover:text-red-300"
        @click="store.clearLogs"
      >
        清除
      </button>
    </div>

    <!-- 日誌內容 -->
    <div ref="logBox" class="flex-1 overflow-y-auto">
      <ul class="text-sm">
        <li
          v-for="(log, idx) in logs"
          :key="idx"
          class="pt-1 px-1"
        >
          <div class="flex items-start text-left w-full bg-stone-900 px-2 py-1 rounded">
            <span class="text-gray-300 w-25">{{ log.timestamp }}</span>

            <span
              class="w-23"
              :class="{
                'text-blue-300'  : log.level === 'info',
                'text-yellow-300': log.level === 'warning',
                'text-red-300'   : log.level === 'error',
              }"
            >
              [{{ log.level }}]
            </span>

            <span class="flex-1 text-left break-normal">{{ log.message }}</span>
          </div>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, watch, nextTick, onMounted } from 'vue'
import { useLoggerStore } from '@/src/store/logger.js'

const store = useLoggerStore()
const logs  = computed(() => store.logs)

const logBox = ref(null)           // outer div for scroll

/* 卷到底部 */
const scrollToBottom = () => {
  if (logBox.value) {
    logBox.value.scrollTop = logBox.value.scrollHeight
  }
}

onMounted(scrollToBottom)

watch(
  logs,
  async () => {
    await nextTick()
    scrollToBottom()
  },
  { deep: true }
)

/* ---------- Props ---------- */
defineProps({
  title: {
    type: String,
    default: '執行日誌',
  },
})
</script>

<style scoped>
ul {
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
}
</style>
