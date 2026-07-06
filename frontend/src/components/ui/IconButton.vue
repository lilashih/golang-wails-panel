<script setup lang="ts">
defineProps<{
  label: string;
}>();
</script>

<template>
  <button
    :aria-label="label"
    :title="label"
    class="inline-flex h-8 w-8 items-center justify-center rounded-md border border-app-border bg-app-surface text-app-text-muted transition hover:bg-app-surface-muted hover:text-app-text focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-app-primary/40 disabled:cursor-not-allowed disabled:opacity-50"
  >
    <slot />
  </button>
</template>
