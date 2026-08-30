<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import { Candy, PanelLeftClose, PanelLeftOpen, PanelRightClose, PanelRightOpen } from "lucide-vue-next";
import { useAppInfoStore } from "@/stores/app_info";
import SidebarNav from "@/components/layout/SidebarNav.vue";
import LiveLogPanel from "@/components/log/LiveLogPanel.vue";
import IconButton from "@/components/ui/IconButton.vue";

const route = useRoute();
const appInfoStore = useAppInfoStore();
const sidebarCollapsed = ref(false);
const logPanelOpen = ref(true);
const appTitle = computed(() => appInfoStore.info.name);
const appVersion = computed(() => appInfoStore.info.version);
const appMode = computed(() => appInfoStore.info.mode === "release" ? "" : appInfoStore.info.mode);

onMounted(() => {
  appInfoStore.loadAppInfo().catch(console.error);
});
</script>

<template>
  <div class="grid h-screen grid-cols-[auto_minmax(0,1fr)_auto] overflow-hidden bg-app-bg text-app-text">
    <aside
      class="flex min-h-0 flex-col border-r border-app-border bg-app-surface transition-[width] duration-200"
      :class="sidebarCollapsed ? 'w-16' : 'w-60'"
    >
      <div class="flex h-14 flex-none items-center justify-between border-b border-app-border px-3">
        <div v-if="!sidebarCollapsed" class="flex min-w-0 items-center gap-3 overflow-hidden">
          <div class="flex h-8 w-8 items-center justify-center rounded-md bg-app-inverse text-sm font-bold text-app-text-inverse">
            <Candy class="h-5 w-5" />
          </div>
          <div class="min-w-0">
            <p class="truncate text-sm font-semibold text-app-text">{{ appTitle }}</p>
            <p class="truncate text-xs text-app-text-muted">版本: {{ appVersion }} {{ appMode }}</p>
          </div>
        </div>
        <IconButton :label="sidebarCollapsed ? '展開左側選單' : '收合左側選單'" @click="sidebarCollapsed = !sidebarCollapsed">
          <PanelLeftOpen v-if="sidebarCollapsed" class="h-4 w-4" />
          <PanelLeftClose v-else class="h-4 w-4" />
        </IconButton>
      </div>
      <SidebarNav class="min-h-0 flex-1" :collapsed="sidebarCollapsed" />
    </aside>

    <section class="flex min-w-0 min-h-0 flex-col overflow-hidden">
      <div class="flex h-14 flex-none items-center justify-between border-b border-app-border bg-app-surface px-4">
        <div class="min-w-0">
          <p class="truncate text-base font-semibold text-app-text text-lg">{{ route.meta.title || ' ' }}</p>
        </div>
        <IconButton :label="logPanelOpen ? '收合右側 log' : '展開右側 log'" @click="logPanelOpen = !logPanelOpen">
          <PanelRightClose v-if="logPanelOpen" class="h-4 w-4" />
          <PanelRightOpen v-else class="h-4 w-4" />
        </IconButton>
      </div>

      <main class="min-h-0 flex-1 overflow-hidden">
        <RouterView v-slot="{ Component }">
          <component :is="Component" class="h-full min-h-0" />
        </RouterView>
      </main>
    </section>

    <section
      class="flex min-h-0 flex-col overflow-hidden border-l border-app-border bg-app-log-bg transition-[width,min-width] duration-200"
      :class="logPanelOpen ? 'w-96 min-w-96' : 'w-0 min-w-0 border-l-0'"
    >
      <LiveLogPanel class="h-full" />
    </section>
  </div>
</template>
