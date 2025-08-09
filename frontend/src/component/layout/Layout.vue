<template>
  <div class="p-6 pb-3">
    <h1 class="text-2xl font-bold mb-5">{{ title }}</h1>

    <ul>
      <li
        v-for="(item, index) in items"
        :key="index"
        class="py-2 px-2"
      >
        <div class="flex items-start text-left w-full bg-stone-900 py-2 pl-4 pr-2 rounded">
          <!--
            將 item 與 index 透過 slot props 傳出去，
            由呼叫端決定要怎麼渲染這一行
          -->
          <slot :item="item" :index="index" />
        </div>
      </li>
    </ul>
  </div>
</template>

<script setup>
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
