<script setup lang="ts">
import { onMounted, ref } from "vue";
import { storeToRefs } from "pinia";
import { Download, FolderOpen, Play, RefreshCw, Square } from "lucide-vue-next";
import { GetAppInfo } from "@bindings/gbase/src/service/app_info/service.js";
import { useProjectsStore } from "@/stores/projects";
import AppButton from "@/components/ui/AppButton.vue";
import EmptyState from "@/components/ui/EmptyState.vue";
import IconButton from "@/components/ui/IconButton.vue";
import StatusBadge from "@/components/ui/StatusBadge.vue";

const projectsStore = useProjectsStore();
const { projects, loaded, error, busyKey, runningCount } = storeToRefs(projectsStore);
const projectPath = ref("");

function displayProjectPath(path: string) {
  const basePath = projectPath.value.trim();
  if (!basePath) {
    return path;
  }

  const normalizedPath = path.replaceAll("\\", "/");
  const normalizedBase = basePath.replaceAll("\\", "/").replace(/\/+$/, "");
  if (normalizedPath.toLowerCase() === normalizedBase.toLowerCase()) {
    return ".";
  }

  const prefix = `${normalizedBase}/`.toLowerCase();
  if (!normalizedPath.toLowerCase().startsWith(prefix)) {
    return path;
  }

  return normalizedPath.slice(normalizedBase.length + 1);
}

onMounted(() => {
  projectsStore.loadProjects().catch(console.error);
  GetAppInfo()
    .then((info) => {
      projectPath.value = `${info.project_path || ""}`.trim();
    })
    .catch(console.error);
});
</script>

<template>
  <section class="flex h-full min-h-0 flex-col overflow-hidden">
    <div class="flex-none border-b border-app-border bg-app-surface px-4 py-3">
      <div class="flex items-center justify-between gap-3">
        <div class="min-w-0">
          <p class="text-sm text-app-text-muted">
            {{ projects.length }} 個項目，{{ runningCount }} 個執行中
          </p>
          <p v-if="projectPath" class="mt-1 truncate font-mono text-xs text-app-text-muted" :title="projectPath">
            {{ projectPath }}
          </p>
        </div>
        <AppButton size="sm" @click="projectsStore.loadProjects">
          <RefreshCw class="h-4 w-4" />
          <span class="text-sm">重新整理</span>
        </AppButton>
      </div>
      <div v-if="error" class="mt-3 rounded-md border border-app-danger/30 bg-app-danger-soft px-3 py-2 text-sm text-app-danger">
        {{ error }}
      </div>
    </div>

    <div class="min-h-0 flex-1 overflow-auto p-4">
      <EmptyState
        v-if="loaded && projects.length === 0"
        title="尚未找到可用專案"
        message="請確認 projects 目錄與 project.json 設定是否正確。"
      />

      <div v-else class="overflow-x-auto rounded-md border border-app-border bg-app-surface">
        <div class="min-w-[40rem]">
          <div class="grid grid-cols-[minmax(3rem,1fr)_6rem_6rem_minmax(5rem,1fr)_10rem] border-b border-app-border bg-app-surface-muted px-4 py-2 text-xs font-semibold uppercase text-app-text-muted text-center">
            <span>名稱</span>
            <span>類型</span>
            <span>狀態</span>
            <span>路徑</span>
            <span>操作</span>
          </div>

          <div
            v-for="project in projects"
            :key="`${project.Config.key}-${project.path}`"
            class="grid min-h-16 grid-cols-[minmax(3rem,1fr)_6rem_6rem_minmax(5rem,1fr)_10rem] items-center border-b border-app-border-muted px-4 py-3 text-center last:border-b-0"
          >
            <div class="min-w-0 text-left">
              <p class="truncate text-sm font-semibold text-app-text">{{ project.Config.title }}</p>
              <p class="mt-0.5 truncate font-mono text-xs text-app-text-muted">{{ project.Config.key }}</p>
            </div>
            <div class="truncate text-sm text-app-text-muted">{{ project.Config.type }}</div>
            <StatusBadge :active="project.running" />
            <div class="truncate px-3 font-mono text-xs text-app-text-muted" :title="project.path">
              {{ displayProjectPath(project.path) }}
            </div>
            <div class="flex justify-end gap-2 text-sm">
              <AppButton
                v-if="!project.running"
                size="sm"
                variant="primary"
                :disabled="busyKey === project.Config.key"
                @click="projectsStore.startProject(project.Config.key)"
              >
                <Play class="h-4 w-4" />
                啟動
              </AppButton>
              <AppButton
                v-else
                size="sm"
                variant="danger"
                :disabled="busyKey === project.Config.key"
                @click="projectsStore.stopProject(project.Config.key)"
              >
                <Square class="h-4 w-4" />
                停止
              </AppButton>
              <IconButton label="重新安裝" :disabled="busyKey === project.Config.key || project.running" @click="projectsStore.installProject(project.Config.key)">
                <Download class="h-4 w-4" />
              </IconButton>
              <IconButton label="開啟資料夾" :disabled="busyKey === project.Config.key" @click="projectsStore.openProject(project.Config.key)">
                <FolderOpen class="h-4 w-4" />
              </IconButton>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
