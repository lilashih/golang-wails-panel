<script setup lang="ts">
withDefaults(
  defineProps<{
    variant?: "primary" | "secondary" | "danger" | "ghost";
    size?: "sm" | "md";
  }>(),
  {
    variant: "secondary",
    size: "md",
  },
);
</script>

<template>
  <button
    :class="[
      'inline-flex items-center justify-center gap-2 rounded-md border font-medium transition focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-app-primary/40 disabled:cursor-not-allowed disabled:opacity-50',
      size === 'sm' ? 'h-8 px-2.5 text-xs' : 'h-9 px-3 text-sm',
      variant === 'primary' && 'border-app-success bg-app-success text-app-text-inverse hover:bg-app-success-hover',
      variant === 'secondary' && 'border-app-border bg-app-surface text-app-text hover:bg-app-surface-muted',
      variant === 'danger' && 'border-app-danger bg-app-danger text-app-text-inverse hover:bg-app-danger-hover',
      variant === 'ghost' && 'border-transparent bg-transparent text-app-text-muted hover:bg-app-bg hover:text-app-text',
    ]"
  >
    <slot />
  </button>
</template>
