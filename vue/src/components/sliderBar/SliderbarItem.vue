<template>
<!-- 递归菜单项组件 -->
<!-- 如果菜单项有子菜单 -->
<el-sub-menu :index="item.path" v-if="item.children">
  <!-- 标题插槽 -->
  <template #title>
    <!-- 菜单图标 -->
    <el-icon><icon-menu /></el-icon>
    <!-- 菜单名称 -->
    <span>{{item.name}}</span>
  </template>
  
  <!-- 递归渲染子菜单项 -->
  <sliderbar-item v-for="(mi,i) in item.children" :index="item.path" :key="mi.id" :item="mi"></sliderbar-item>
</el-sub-menu>
  
<!-- 如果没有子菜单 -->
<el-menu-item :index="item.path" v-else>
  <!-- 菜单图标 -->
  <el-icon><icon-menu /></el-icon>
  <!-- 菜单标题 -->
  <template #title>{{item.name}}</template>
</el-menu-item>
</template>

<script lang="ts" setup>
// 导入Element Plus图标
import {
  Location,
  Document,
  Menu as IconMenu,
  Setting,
} from '@element-plus/icons-vue'

// 定义组件属性
const props=defineProps({
    // 菜单项数据
    item:{
        type:Object,
        default(){
            return {}
        }
    },
    // 菜单项索引
    index:{}
})
</script>

<style scoped>
/* 空样式（保留原始结构） */
</style>