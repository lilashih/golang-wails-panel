<script setup lang="ts">
import { computed, nextTick, ref, watch } from "vue";
import { storeToRefs } from "pinia";
import { ChevronsDown, Eraser } from "lucide-vue-next";
import { useLogsStore } from "@/stores/logs";
import LogLevelBadge from "@/components/log/LogLevelBadge.vue";

const logsStore = useLogsStore();
const { logs } = storeToRefs(logsStore);
const panelRef = ref<HTMLElement | null>(null);
const follow = ref(true);

const total = computed(() => logs.value.length);

function scrollToBottom() {
  nextTick(() => {
    const el = panelRef.value;
    if (!el) {
      return;
    }
    el.scrollTo({ top: el.scrollHeight, behavior: "auto" });
  });
}

watch(
  () => logs.value.length,
  () => {
    if (follow.value) {
      scrollToBottom();
    }
  },
  { flush: "post" },
);

watch(follow, (enabled) => {
  if (enabled) {
    scrollToBottom();
  }
});
</script>

<template>
  <aside class="flex h-full min-h-0 flex-col overflow-hidden bg-app-log-bg text-app-text-inverse">
    <header class="flex h-14 flex-none items-center justify-between border-b border-app-log-divider/10 px-3">
      <div>
        <p class="text-sm font-semibold">即時 Log</p>
        <p class="text-xs text-app-text-subtle">{{ total }} 筆</p>
      </div>
      <div class="flex items-center gap-2">
        <label class="flex h-8 items-center gap-2 px-2 text-xs text-app-text-inverse/80">
          <input v-model="follow" type="checkbox" class="h-3.5 w-3.5" />
          強制置底
        </label>
        <div class="flex h-8 items-center rounded-md border border-app-log-divider/10 bg-app-log-divider/[0.03] p-0.5">
          <button
            type="button"
            aria-label="跳到底部"
            title="跳到底部"
            class="inline-flex h-7 w-7 items-center justify-center rounded text-app-text-subtle transition hover:text-app-text-inverse focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-app-primary/40"
            @click="scrollToBottom"
          >
            <ChevronsDown class="h-4 w-4" />
          </button>
          <button
            type="button"
            aria-label="清除"
            title="清除"
            class="inline-flex h-7 w-7 items-center justify-center rounded text-app-text-subtle transition hover:text-app-text-inverse focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-app-primary/40"
            @click="logsStore.clearLogs()"
          >
            <Eraser class="h-4 w-4" />
          </button>
        </div>
      </div>
    </header>

    <div ref="panelRef" class="min-h-0 flex-1 overflow-auto px-3 py-2 font-mono text-xs leading-5">
      <div v-if="logs.length === 0" class="flex h-full items-center justify-center text-app-text-muted">
        尚無即時訊息
      </div>
      <div
        v-for="(item, index) in logs"
        :key="`${item.timestamp}-${index}`"
        class="grid grid-cols-[3.5rem_2.8rem_1fr] gap-2 border-b border-app-log-divider/5 py-1"
      >
        <span class="mono-tabular text-app-text-muted">{{ item.timestamp }}</span>
        <LogLevelBadge :level="item.level" />
        <span class="break-words text-app-text-inverse/90">{{ item.message }}</span>
      </div>
    </div>
  </aside>
</template>
