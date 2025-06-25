<template>
<!-- 根容器，根据collapse状态切换类名 -->
<div :class="['slide-bar',collapse?'collapse':'']">
    <!-- 未折叠时显示的Logo -->
    <div class="logo" v-if="!collapse">
       <!-- 路由链接到首页 -->
       <router-link to="/index">在线练题系统</router-link>
     </div>
     <!-- 折叠时显示的简化Logo -->
     <div class="logo" v-else>
       <router-link to="/index">logo</router-link>
     </div>
     <!-- 导航栏容器 -->
     <div class="headers">
 <!-- Element Plus菜单组件 -->
 <el-menu
        active-text-color="#dfff7d"
        background-color="#0087bf"
        class="el-menu-vertical-demo"
         :default-active="onRoutes" 
        text-color="#fff"
        style="border:none"
        :collapse="collapse"
        router
         mode="horizontal"
      >
       <!-- 循环渲染菜单项 -->
       <sliderbar-item v-for="(item,index) in menuList" :index="index" :key="item.id" :item="item"></sliderbar-item>  
      </el-menu>
       <!-- 引入头部组件 -->
       <v-header></v-header>
     </div>
</div>
</template>

<script lang="ts" setup>
// 导入Element Plus图标
import {
  Location,
  Document,
  Menu as IconMenu,
  Setting,
} from '@element-plus/icons-vue'
// 导入Vue组合式API
import {computed, ref} from 'vue'
// 导入Vuex和Vue Router
import {useStore} from 'vuex'
import {useRoute} from 'vue-router'
// 导入自定义组件
import SliderbarItem from './sliderBar/SliderbarItem.vue'
import vHeader from "../components/Header.vue";

// 获取Vuex store实例
const store=useStore()
// 计算折叠状态（从Vuex获取）
const collapse=computed(()=>store.state.collapse)

// 获取当前路由信息
const route=useRoute()
// 计算当前激活的路由路径
const onRoutes = computed(() => {
            return route.path;
        });

// 菜单数据
const menuList=[
  {name:'问题列表',id:1,path:'/questionList' },
  {name:'提交列表',id:2,path:'/submitList'},
  {name:'排行榜',id:3,path:'/topList'},
]
</script>

<style>
/* 全局样式（空） */
</style>

<style scoped lang="scss">
/* 带作用域的SCSS样式 */

/* Logo容器样式 */
.logo{
  height: 50px;
  display: none;  /* 默认隐藏 */
  background-color: #0087bf;
  border-bottom: 1px solid #dfff7d;
  align-items: center;
  justify-content: center;
  /* 链接样式 */
  a{
    color: #fff;
    font-weight: 600;
  }
}

/* 侧边栏容器样式 */
.slide-bar{
     transition: 0.3s;  /* 过渡动画 */
     min-width: 250px;
     background-color: #0087bf;
     /* 折叠状态样式 */
     &.collapse{
       min-width: 0;
     }
}

/* Element菜单容器 */
.el-menu-vertical-demo{
     flex: 1;  /* 弹性填充 */
}

/* 头部容器样式 */
.headers{
     display: flex;
     align-items: center;
}
</style>