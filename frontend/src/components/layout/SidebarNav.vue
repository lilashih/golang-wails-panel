<script setup lang="ts">
import { computed, onMounted, ref, type Component } from "vue";
import { RouterLink } from "vue-router";
import { GetAppInfo } from "@bindings/gbase/src/service/app_info/service.js";
import { routes } from "@/router/routes";

const props = defineProps<{
  collapsed: boolean;
}>();

const appMode = ref("release");

const navItems = computed(() =>
  routes
    .filter((route) => route.meta?.isMenu)
    .filter((route) => !route.meta?.devOnly || appMode.value !== "release")
    .map((route) => ({
      to: route.path,
      label: String(route.meta?.title || route.name || route.path),
      icon: route.meta?.icon as Component,
    })),
);

onMounted(async () => {
  try {
    const info = await GetAppInfo();
    appMode.value = `${info.mode || "release"}`.trim() || "release";
  } catch (err) {
    console.error(err);
  }
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
