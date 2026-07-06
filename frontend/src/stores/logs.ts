import { computed, ref } from "vue";
import { defineStore } from "pinia";
import type { LiveLogEvent } from "@/types/log";

function formatLogTime(date = new Date()) {
  return date.toLocaleTimeString("en-GB", {
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
    hour12: false,
  });
}

export const useLogsStore = defineStore("logs", () => {
  const logs = ref<LiveLogEvent[]>([]);
  const maxItems = ref(500);
  const unreadCount = ref(0);

  const latest = computed(() => logs.value[logs.value.length - 1] ?? null);

  function addLog(payload: LiveLogEvent) {
    logs.value.push({
      level: payload.level || "info",
      message: payload.message || "",
      timestamp: payload.timestamp || formatLogTime(),
    });

    if (logs.value.length > maxItems.value) {
      logs.value.splice(0, logs.value.length - maxItems.value);
    }
    unreadCount.value += 1;
  }

  function clearLogs() {
    logs.value = [];
    unreadCount.value = 0;
  }

  function markRead() {
    unreadCount.value = 0;
  }

  return {
    logs,
    maxItems,
    unreadCount,
    latest,
    addLog,
    clearLogs,
    markRead,
  };
});
