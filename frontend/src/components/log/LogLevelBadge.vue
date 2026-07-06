<script setup lang="ts">
const props = defineProps<{
  level: string;
}>();

function levelClass() {
  switch (props.level) {
    case "error":
      return "bg-app-danger-soft text-app-danger ring-app-danger/30";
    case "warning":
      return "bg-app-warning-soft text-app-warning ring-app-warning/30";
    case "debug":
      return "bg-app-bg text-app-text-muted ring-app-border";
    default:
      return "bg-app-primary-soft text-app-primary ring-app-primary/30";
  }
}
</script>

<template>
  <span :class="['inline-flex w-11 h-6 justify-center rounded px-1.5 py-0.5 text-[11px] font-semibold uppercase ring-1', levelClass()]">
    {{ level || "info" }}
  </span>
</template>
