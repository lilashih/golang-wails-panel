<template>
<div class="flex flex-col flex-1 overflow-hidden">
  <div class="flex-1 overflow-auto p-6">
    <Layout :title="title" :items="items">
      <template #default="{ item, index }">
        <slot :item="item" :index="index" />
      </template>
    </Layout>
  </div>
  <div class="shadow-inner h-2/5 p-3">
    <Logger />
  </div>
</div>
</template>

<script setup>
import Layout from '@/src/component/layout/Layout.vue'
import Logger from '@/src/component/widget/Logger.vue'

/* ---------- Props ---------- */
defineProps({
  title: {
    type: String,
    default: '標題',
  },
  items: {
    type: Array,
    default: () => [],
  },
})
</script>
