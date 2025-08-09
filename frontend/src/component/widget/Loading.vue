<template>
  <div v-if="visible" class="loading-mask">
    <div class="loading-content">
      <!-- 旋轉 + 跑描邊 的圓形 Loader -->
      <svg class="spinner" viewBox="0 0 50 50">
        <circle class="path" cx="25" cy="25" r="20" fill="none" />
      </svg>
      <div class="loading-text">{{ text }}</div>
    </div>
  </div>
</template>

<script setup>
/* ---------- Props ---------- */
defineProps({
  text: {
    type: String,
    default: '',
  },
  visible: {
    type: Boolean,
    default: false,
  },
})
</script>

<style scoped>
/* 遮罩 */
.loading-mask {
  position: fixed;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(3, 3, 3, 0.4);
  z-index: 9999;
}

/* 內容容器 */
.loading-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.75rem;
}

/* SVG 旋轉動畫 */
.spinner {
  width: 56px;
  height: 56px;
  animation: spinner-rotate 2s linear infinite;
}

/* 圓形描邊動畫 */
.path {
  stroke: #636e7b;
  stroke-width: 4;
  stroke-linecap: round;
  /* dasharray / dashoffset 搭配 keyframes 產生“追光”效果 */
  stroke-dasharray: 150, 200;
  stroke-dashoffset: -10;
  animation: spinner-dash 1.5s ease-in-out infinite;
}

/* 文字 */
.loading-text {
  font-size: 1.125rem; /* 18px */
  color: #fff;
  letter-spacing: 2px;
}

/* ───────── Keyframes ───────── */
@keyframes spinner-rotate {
  100% { transform: rotate(360deg); }
}

@keyframes spinner-dash {
  0%   { stroke-dasharray: 1, 200;  stroke-dashoffset: 0;   }
  50%  { stroke-dasharray: 90, 200; stroke-dashoffset: -35; }
  100% { stroke-dasharray: 90, 200; stroke-dashoffset: -124;}
}
</style>
