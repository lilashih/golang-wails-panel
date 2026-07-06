import type { RouteRecordRaw } from "vue-router";
import { BookOpen, FolderKanban, ScrollText } from "lucide-vue-next";
import PanelPage from "@/pages/PanelPage.vue";
import LogViewerPage from "@/pages/LogViewerPage.vue";
import DevGuidePage from "@/pages/DevGuidePage.vue";

export const routes: RouteRecordRaw[] = [
  { path: "/", redirect: "/panel" },
  {
    path: "/panel",
    name: "Panel",
    component: PanelPage,
    meta: { title: "應用程式", icon: FolderKanban, isMenu: true },
  },
  {
    path: "/log-viewer",
    name: "LogViewer",
    component: LogViewerPage,
    meta: { title: "日誌檢視", icon: ScrollText, isMenu: true },
  },
  {
    path: "/dev-guide",
    name: "DevGuide",
    component: DevGuidePage,
    meta: { title: "開發指南", icon: BookOpen, isMenu: true, devOnly: true },
  },
];
