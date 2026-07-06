<script setup lang="ts">
defineProps<{
  title: string;
  message?: string;
}>();
</script>

<template>
  <div class="flex min-h-48 flex-col items-center justify-center rounded-md border border-dashed border-app-border bg-app-surface/70 px-6 py-10 text-center">
    <p class="text-sm font-semibold text-app-text">{{ title }}</p>
    <p v-if="message" class="mt-1 max-w-md text-sm text-app-text-muted">{{ message }}</p>
  </div>
</template>
