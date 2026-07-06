<script setup lang="ts">
import { storeToRefs } from "pinia";
import { LoaderCircle } from "lucide-vue-next";
import { useLoadingStore } from "@/stores/loading";

const loadingStore = useLoadingStore();
const { loading, message } = storeToRefs(loadingStore);
</script>

<template>
  <div
    v-if="loading"
    class="fixed inset-0 z-50 flex items-center justify-center bg-app-overlay/20 backdrop-blur-[2px]"
  >
    <div class="flex items-center gap-3 rounded-md border border-app-border bg-app-surface px-4 py-3 text-sm font-medium text-app-text shadow-lg">
      <LoaderCircle class="h-4 w-4 animate-spin text-app-primary" />
      <span>{{ message || "處理中" }}</span>
    </div>
  </div>
</template>
