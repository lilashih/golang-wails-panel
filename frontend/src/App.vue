<template>
  <div class="flex h-screen">
    <!-- 左側選單 -->
    <nav class="w-56 h-screen bg-neutral-900 flex flex-col pt-8">
      <ul class="space-y-2">
        <li v-for="item in pages" :key="item.path">
          <RouterLink class="menu" active-class="menu-active" :to="item.path">
            {{ item.label }}
          </RouterLink>
        </li>
      </ul>
    </nav>

    <!-- 右側內容 -->
    <div class="flex flex-col flex-1 h-full">
      <RouterView />
      <Loading :visible="loading" :text="loadingText" />
    </div>
  </div>
</template>

<script setup>
import Loading from '@/src/component/widget/Loading.vue'
import { storeToRefs } from 'pinia'
import { useLoadingStore } from '@/src/store/loading.js'
import routes from '@/src/router/page.js'

const pages = routes.map(r => ({
  path: r.path,
  label: r.title || r.name || r.path
}))

const loadingStore = useLoadingStore()
const { loading, loadingText } = storeToRefs(loadingStore)
</script>

<style scoped>
@reference "@/src/css/style.css";

.menu {
  @apply block px-2 m-3 py-2 text-gray-300;
}
.menu:hover,
.menu-active {
  @apply bg-stone-700 text-white rounded-lg;
}
</style>