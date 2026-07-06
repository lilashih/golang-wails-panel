import { computed, ref } from "vue";
import { defineStore } from "pinia";
import { useLoadingStore } from "@/stores/loading";
import { useLogsStore } from "@/stores/logs";
import {
  InstallProject,
  ListProjects,
  OpenProject,
  StartProject,
  StopProject,
} from "@bindings/gbase/src/service/panel/service.js";

export interface ProjectConfig {
  os: string;
  title: string;
  key: string;
  type: string;
  start: string;
  stop: string;
  install: string;
}

export interface ProjectItem {
  Config: ProjectConfig;
  path: string;
  running: boolean;
}

function toErrorMessage(err: unknown) {
  return err instanceof Error ? err.message : String(err);
}

function isProjectConfig(value: unknown): value is ProjectConfig {
  if (!value || typeof value !== "object") {
    return false;
  }

  const config = value as Partial<ProjectConfig>;
  return [
    config.os,
    config.title,
    config.key,
    config.type,
    config.start,
    config.stop,
    config.install,
  ].every((item) => typeof item === "string");
}

function isProjectItem(value: unknown): value is ProjectItem {
  if (!value || typeof value !== "object") {
    return false;
  }

  const project = value as Partial<ProjectItem>;
  return isProjectConfig(project.Config) && typeof project.path === "string" && typeof project.running === "boolean";
}

function normalizeProjects(value: unknown): ProjectItem[] {
  if (!Array.isArray(value)) {
    return [];
  }

  return value.filter(isProjectItem);
}

export const useProjectsStore = defineStore("projects", () => {
  const projects = ref<ProjectItem[]>([]);
  const loaded = ref(false);
  const error = ref("");
  const busyKey = ref<string | null>(null);

  const runningCount = computed(() => projects.value.filter((project) => project.running).length);

  async function refreshProjects() {
    const result = await ListProjects();
    projects.value = normalizeProjects(result);
    loaded.value = true;
  }

  async function loadProjects() {
    const loading = useLoadingStore();
    error.value = "";
    try {
      await loading.withLoading("讀取中", refreshProjects);
    } catch (err) {
      error.value = toErrorMessage(err);
      loaded.value = true;
    }
  }

  async function runProjectAction(key: string, label: string, task: () => Promise<unknown>) {
    const loading = useLoadingStore();
    const logs = useLogsStore();
    busyKey.value = key;
    error.value = "";
    try {
      await loading.withLoading(`${label}中`, async () => {
        await task();
        await refreshProjects();
      });
    } catch (err) {
      const message = toErrorMessage(err);
      error.value = message;
      logs.addLog({
        level: "error",
        message: `${label}失敗：${message}`,
        timestamp: "",
      });
    } finally {
      busyKey.value = null;
    }
  }

  function openProject(key: string) {
    return runProjectAction(key, "開啟資料夾", () => OpenProject(key));
  }

  function startProject(key: string) {
    return runProjectAction(key, "啟動", () => StartProject(key));
  }

  function stopProject(key: string) {
    return runProjectAction(key, "停止", () => StopProject(key));
  }

  function installProject(key: string) {
    return runProjectAction(key, "重新安裝", () => InstallProject(key));
  }

  return {
    projects,
    loaded,
    error,
    busyKey,
    runningCount,
    loadProjects,
    openProject,
    startProject,
    stopProject,
    installProject,
  };
});
