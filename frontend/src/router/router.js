import { createRouter, createWebHashHistory } from 'vue-router'
import { useLoggerStore } from '@/src/store/logger.js'
import { useLoadingStore } from '@/src/store/loading.js'
import pages from '@/src/router/page.js';

export const router = createRouter({
  // Wails 桌面環境建議用 Hash；重新整理不會 404
  history: createWebHashHistory(),
  routes: [
    { path: '/', redirect: '/panel' },
    ...pages,
  ],
})

// 切換頁面時顯示 loading 遮罩
router.beforeEach((to, from, next) => {
  const loadingStore = useLoadingStore()
  loadingStore.setLoading(true)
  next()
})

router.afterEach(() => {
  // 必須在回呼裡面拿 Store，才能取得目前 App 內的實例

  // 路由切換後立即清空 log
  // const loggerStore = useLoggerStore()
  // loggerStore.clearLogs()

  // 關閉 loading 遮罩
  const loadingStore = useLoadingStore()
  loadingStore.setLoading(false)
})