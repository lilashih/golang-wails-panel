import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useLoadingStore = defineStore('loading', () => {
  const loading = ref(false)
  const loadingText = ref('')

  const setLoading = (val, text = '') => {
    loading.value = val
    loadingText.value = text
  }

  return {
    loading,
    loadingText,
    setLoading,
  }
})
