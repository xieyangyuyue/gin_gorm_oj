<template>
  <!-- 主页容器 -->
  <div>
    <!-- 内容区域（当前为空） -->
    <!-- 根据业务需求添加内容 -->
  </div>
</template>

<script lang="ts" setup>
// 导入 Element Plus 图标组件
import { DArrowRight } from '@element-plus/icons-vue'

// Vue 响应式函数
import { ref, onMounted } from "vue";

// 按需引入 ECharts 图表库
import * as echarts from "echarts";

// 图表容器引用数组
const arr = ref([])

// 测试数据（图表数据）
const brr = [
  { name: '已通知的考试人数', value: 12, color: '#16CEB9' },
  { name: '已完成考试人数', value: 23, color: '#6648FF' },
  { name: '未完成考试人数', value: 14, color: '#399' },
  { name: '已超期考试人数', value: 32, color: '#399' }
]

/**
 * 图表初始化函数
 * @param {HTMLElement} el - 图表挂载的DOM元素
 */
const charts = (el: HTMLElement) => {
  // 将图表容器添加到引用数组
  (arr.value as Array<HTMLElement>).push(el)
}

// 图表配置项接口（定义数据结构）
interface chartOption {
  el: HTMLElement;    // 图表挂载元素
  name: string;       // 图表名称
  num: number,        // 数值
  color: string;      // 颜色
}
</script>

<style scoped lang="scss">
/* 卡片样式 */
.card {
  padding: 10px;        /* 内边距 */
  box-sizing: border-box; /* 包含padding在尺寸内 */
}

/* 顶部栏样式 */
.bar {
  display: flex;             /* 弹性布局 */
  justify-content: space-between; /* 两端对齐 */
  background-color: #f5f5f5; /* 浅灰色背景 */
  line-height: 30px;         /* 行高 */
  align-items: center;       /* 垂直居中 */
  padding: 10px 20px;        /* 内边距 */
  
  /* 第二个span元素样式 */
  span:nth-child(2) {
    color: #999;      /* 浅灰色文字 */
    font-size: 14px; /* 较小字体 */
  }
}

/* 图表容器样式 */
.charts {
  display: flex;             /* 弹性布局 */
  justify-content: space-between; /* 项目均匀分布 */
  height: 200px;             /* 固定高度 */
  
  /* 直接子元素样式 */
  &>div {
    width: 25%; /* 每个图表占25%宽度 */
  }
}

/* 标题样式 */
h3 {
  border-bottom: 1px solid #eee; /* 底部边框 */
  padding: 20px;                 /* 内边距 */
  margin-bottom: 10px;           /* 底部外边距 */
}
</style>