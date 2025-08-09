import { defineConfig } from 'vite'
import { resolve } from 'path'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'

// https://vitejs.dev/config/
export default defineConfig({
  root: resolve(__dirname),
  resolve: {
    // 使用 @ 符號作為前端根目錄路徑別名
    alias: [
      {
        find: '@',
        replacement: resolve(__dirname),
      },
    ],
  },
  plugins: [
    vue(),
    tailwindcss(),
  ],
})
