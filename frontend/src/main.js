import {createApp} from 'vue'
import { createPinia } from 'pinia'
import VueVirtualScroller from 'vue3-virtual-scroller'
import 'vue3-virtual-scroller/dist/vue3-virtual-scroller.css'
import App from './App.vue'
import '@/src/css/style.css';
import { EventsOn } from '@/wailsjs/runtime'
import { useLoggerStore } from '@/src/store/logger.js'
import { router } from '@/src/router/router.js';

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(VueVirtualScroller)

const loggerStore = useLoggerStore()
EventsOn('log', (payload) => {
  loggerStore.addLog(payload.level, payload.message)
})

app.mount('#app')