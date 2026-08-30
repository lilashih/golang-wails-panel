<script setup lang="ts">
import { computed, onMounted, type Component } from "vue";
import { RouterLink } from "vue-router";
import { routes } from "@/router/routes";
import { useAppInfoStore } from "@/stores/app_info";

const props = defineProps<{
  collapsed: boolean;
}>();

const appInfoStore = useAppInfoStore();

const navItems = computed(() =>
  routes
    .filter((route) => route.meta?.isMenu)
    .filter((route) => !route.meta?.devOnly || appInfoStore.info.mode !== "release")
    .map((route) => ({
      to: route.path,
      label: String(route.meta?.title || route.name || route.path),
      icon: route.meta?.icon as Component,
    })),
);

onMounted(() => {
  appInfoStore.loadAppInfo().catch(console.error);
});
</script>

<template>
  <nav class="min-h-0 overflow-auto px-2 py-3">
    <RouterLink
      v-for="item in navItems"
      :key="item.to"
      :to="item.to"
      custom
      v-slot="{ href, navigate, isActive }"
    >
      <a
        :href="href"
        :title="item.label"
        @click="navigate"
        class="flex h-10 items-center gap-3 rounded-md text-sm font-medium transition"
        :class="[
          collapsed ? 'justify-center px-0' : 'px-3',
          isActive ? 'bg-app-inverse text-app-text-inverse hover:bg-app-inverse hover:text-app-text-inverse' : 'text-app-text-muted hover:bg-app-bg hover:text-app-text',
        ]"
      >
        <component :is="item.icon" class="h-4 w-4 shrink-0" />
        <span v-if="!collapsed" class="truncate">{{ item.label }}</span>
      </a>
    </RouterLink>
  </nav>
</template>
