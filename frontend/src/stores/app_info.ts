import { ref } from "vue";
import { defineStore } from "pinia";
import { GetAppInfo } from "@bindings/gbase/src/service/app_info/service.js";
import type { AppInfo } from "@bindings/gbase/src/service/app_info/models.js";

const defaultAppInfo: AppInfo = {
  name: "GBase",
  version: "0.0.0",
  mode: "release",
  project_path: "",
};

function normalizeAppInfo(value: AppInfo): AppInfo {
  return {
    name: `${value.name || defaultAppInfo.name}`.trim() || defaultAppInfo.name,
    version: `${value.version || defaultAppInfo.version}`.trim() || defaultAppInfo.version,
    mode: `${value.mode || defaultAppInfo.mode}`.trim() || defaultAppInfo.mode,
    project_path: `${value.project_path || ""}`.trim(),
  };
}

export const useAppInfoStore = defineStore("appInfo", () => {
  const info = ref<AppInfo>({ ...defaultAppInfo });
  const loaded = ref(false);
  const loading = ref(false);
  const error = ref("");
  let loadingPromise: Promise<AppInfo> | null = null;

  function loadAppInfo() {
    if (loaded.value) {
      return Promise.resolve(info.value);
    }

    if (loadingPromise) {
      return loadingPromise;
    }

    loading.value = true;
    error.value = "";
    loadingPromise = GetAppInfo()
      .then((result) => {
        info.value = normalizeAppInfo(result);
        loaded.value = true;
        return info.value;
      })
      .catch((err: unknown) => {
        error.value = err instanceof Error ? err.message : String(err);
        throw err;
      })
      .finally(() => {
        loading.value = false;
        loadingPromise = null;
      });

    return loadingPromise;
  }

  return {
    info,
    loaded,
    loading,
    error,
    loadAppInfo,
  };
});
