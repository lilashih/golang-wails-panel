import { ref } from "vue";
import { defineStore } from "pinia";

interface ActiveLoading {
  id: number;
  message: string;
}

export const useLoadingStore = defineStore("loading", () => {
  const loading = ref(false);
  const message = ref("");
  const pendingCount = ref(0);
  const activeLoadings = ref<ActiveLoading[]>([]);
  let nextLoadingId = 0;
  let manualLoadingId: number | null = null;

  function syncLoadingState() {
    pendingCount.value = activeLoadings.value.length;
    loading.value = pendingCount.value > 0;
    message.value = activeLoadings.value.at(-1)?.message || "";
  }

  function beginLoading(nextMessage: string) {
    const id = nextLoadingId++;
    activeLoadings.value = [...activeLoadings.value, { id, message: nextMessage }];
    syncLoadingState();
    return id;
  }

  function endLoading(id: number | null) {
    if (id === null) {
      return;
    }

    activeLoadings.value = activeLoadings.value.filter((item) => item.id !== id);
    syncLoadingState();
  }

  function setLoading(value: boolean, nextMessage = "") {
    if (value) {
      if (manualLoadingId === null) {
        manualLoadingId = beginLoading(nextMessage);
        return;
      }

      activeLoadings.value = activeLoadings.value.map((item) => (
        item.id === manualLoadingId ? { ...item, message: nextMessage } : item
      ));
      syncLoadingState();
      return;
    }

    endLoading(manualLoadingId);
    manualLoadingId = null;
  }

  async function withLoading<T>(nextMessage: string, task: () => Promise<T>) {
    const id = beginLoading(nextMessage);
    try {
      return await task();
    } finally {
      endLoading(id);
    }
  }

  return {
    loading,
    message,
    pendingCount,
    setLoading,
    withLoading,
  };
});
