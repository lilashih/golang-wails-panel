import Panel from '@/src/component/page/Panel.vue'
import LogViewer from '@/src/component/page/LogViewer.vue'

const routes = [
  { path: '/panel', name: 'Panel', title: '應用程式', component: Panel },
  { path: '/log-viewer', name: 'LogViewer', title: '日誌檢視', component: LogViewer },
]

export default routes;