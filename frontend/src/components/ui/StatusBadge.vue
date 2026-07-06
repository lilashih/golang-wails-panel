<script setup lang="ts">
defineProps<{
  active: boolean;
}>();
</script>

<template>
  <span
    :class="[
      'inline-flex h-6 min-w-10 items-center justify-center rounded-full px-2 m-2 text-xs font-semibold mono-tabular',
      active ? 'bg-app-success-soft text-app-success ring-1 ring-app-success/30' : 'bg-app-bg text-app-text-muted ring-1 ring-app-border',
    ]"
  >
    {{ active ? "執行中" : "已停止" }}
  </span>
</template>
