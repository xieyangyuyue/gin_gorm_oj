<script lang="ts" setup>
// 导入侧边栏组件
import vSlideBar from "../components/SliderBar.vue";
// 导入历史标签组件（可能是面包屑导航或标签页）
import vHistoryTag from "../components/historyTag.vue";
</script>

<template>
  <!-- 页面主容器 -->
  <div class="page">
    <!-- 侧边栏组件 -->
    <v-slide-bar></v-slide-bar>
    
    <!-- 主内容区域 -->
    <div class="content">
      <!-- 顶部导航区域 -->
      <div class="top-bar">
        <!-- 历史标签组件（可能是面包屑或页面标签） -->
        <v-history-tag></v-history-tag>
      </div>

      <!-- 页面内容容器 -->
      <div class="page-box">
        <!-- 页面内容区域 -->
        <div class="page-cont">
          <!-- Vue Router 的视图容器 -->
          <router-view v-slot="{ Component }">
            <!-- 路由切换时的过渡动画 -->
            <transition name="component-fade" mode="out-in">
              <!-- 动态渲染当前路由对应的组件 -->
              <component :is="Component" />
            </transition>
          </router-view>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
/* 路由切换时的淡入淡出动画效果 */
.component-fade-enter-active,
.component-fade-leave-active {
  transition: opacity 0.3s ease; /* 0.3秒透明度过渡效果 */
}

.component-fade-enter-from,
.component-fade-leave-to {
  opacity: 0; /* 进入前和离开后的透明度为0 */
}

/* 页面整体布局样式 */
.page {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column; /* 垂直方向排列 */

  /* 主内容区样式 */
  .content {
    flex: 1; /* 占据剩余空间 */
    overflow: hidden; /* 隐藏溢出内容 */
    
    /* 顶部导航栏样式 */
    .top-bar {
      // height: 80px; 注释掉的高度设置
      box-shadow: 0 10px 10px #eee; /* 底部阴影效果 */
    }
    
    /* 头部区域样式（未使用） */
    .header {
      height: 40px;
      border-bottom: 1px solid #eee;
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
    
    /* 页面内容框样式 */
    .page-box {
      // height: calc(100% - 80px); 注释掉的高度计算
      height: 100%; /* 填充父容器高度 */
      padding: 10px; /* 内边距 */
      box-sizing: border-box; /* 包含padding在尺寸内 */
      background-color: #f5f5f5; /* 浅灰色背景 */
      
      /* 页面内容容器样式 */
      .page-cont {
        box-shadow: 0 0 10px #eee; /* 四周阴影 */
        width: 100%;
        height: 100%;
        overflow-y: auto; /* 垂直滚动 */
        overflow-x: hidden; /* 隐藏水平滚动 */
        background-color: white; /* 白色背景 */
        
        /* 组件容器样式（未使用） */
        .comp-box {
          padding: 10px;
          height: 100%;
          box-sizing: border-box;
        }
      }
    }
  }
}
</style>