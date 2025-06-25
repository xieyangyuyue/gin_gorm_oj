<template>
  <!-- 问题详情容器 -->
  <div class="ques-cont">
    <!-- 左侧详情区域 -->
    <div class="left">
      <!-- 问题标题 -->
      <h3>{{detail.title}}</h3>
      
      <!-- 问题分类信息 -->
      <div class="msg">
        分类：<span v-for="mi in detail.problem_categories" :key="mi.id">
          {{mi.category_basic.name}}  <!-- 显示分类名称 -->
        </span>
      </div>
      
      <!-- 问题元数据 -->
      <div class="msg">
        <span>创建时间{{detail.created_at}}</span> <!-- 创建时间 -->
        <span>提交次数：{{detail.submit_num}}</span> <!-- 提交次数 -->
        <span>通过次数：{{detail.pass_num}}</span> <!-- 通过次数 -->
        <span>最大耗时：{{detail.max_runtime}}ms</span> <!-- 最大运行时间 -->
      </div>
      
      <!-- 问题内容（HTML格式） -->
      <div v-html="detail.content"></div> <!-- 显示问题描述内容 -->
    </div>
    
    <!-- 右侧编辑器区域 -->
    <div class="right">
      <!-- 代码编辑器组件 -->
      <Editor></Editor>  
    </div>
  </div>
</template>

<script lang="ts" setup>
// 导入Vue响应式API
import { reactive, ref } from '@vue/reactivity'

// 导入路由功能
import { useRoute } from 'vue-router'

// 导入API接口
import api from '../../api/api.js'

// 导入代码编辑器组件
import Editor from './editor.vue'

// 导入Element Plus图标（未使用）
import { Edit } from '@element-plus/icons-vue'

// 获取当前路由信息
const route = useRoute()

// 问题详情响应式数据
const detail = ref({}) // 初始化为空对象

// 获取问题详情函数
const getDetail = () => {
  // 调用API获取问题详情
  api.getProblemDetail({
    identity: route.query.identity // 从路由参数获取问题标识
  }).then(res => {
    if (res.data.data) {
      detail.value = res.data.data // 更新问题详情数据
    }
  })
}

// 初始化时获取问题详情
getDetail()
</script>

<style scoped lang="scss">
/* 问题详情整体布局 */
.ques-cont {
  display: flex;    /* 弹性布局 */
  height: 100%;     /* 填充父容器高度 */
  
  /* 左侧详情区域样式 */
  .left {
    width: 50%;          /* 占据50%宽度 */
    line-height: 2;      /* 行高为字体2倍 */
    border-right: 10px solid #eee; /* 右侧分隔线 */
    padding: 10px;       /* 内边距 */
    
    /* 元数据样式 */
    .msg {
      font-size: 12px;   /* 较小字体 */
      
      /* 数据项样式 */
      span {
        margin-right: 20px; /* 右侧外边距 */
        color: #999;        /* 灰色文字 */
      }
    }
  }
  
  /* 右侧编辑器区域样式 */
  .right {
    width: 50%; /* 占据50%宽度 */
  }
}
</style>