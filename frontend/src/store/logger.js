import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useLoggerStore = defineStore('logs', () => {
  const logs = ref([])

  const addLog = (level, message) => {
    logs.value.push({
      timestamp: new Date().toLocaleTimeString(),
      level,
      message
    })
  }

  const info = (message) => addLog('info', message)
  const warning = (message) => addLog('warning', message)
  const error = (message) => addLog('error', message)

  const clearLogs = () => {
    logs.value = []
  }

  return {
    logs,
    addLog,
    info,
    warning,
    error,
    clearLogs
  }
})
