<template>
  <!-- 排行榜容器 -->
  <div class="top-list">
    <!-- 排行榜列表区域 -->
    <div class="list">
      <!-- 遍历排行榜项 -->
      <div class="item" v-for="item in rankList" :key="item.id">
        <!-- 标题区域 -->
        <div class="title">
          {{item.id}}. <!-- 排名序号 -->
          {{item.name}} <!-- 用户名 -->
          <!-- 排序信息 -->
          <div class="sort">
            <span>{{item.mail}}</span> <!-- 用户邮箱 -->
          </div>
        </div>
        <!-- 用户提交信息 -->
        <div class="msg">
          <span>提交数：{{item.submit_num}}</span> <!-- 提交次数 -->
          <span>通过数：{{item.pass_num}}</span> <!-- 通过次数 -->
        </div>
        
        <!-- 编辑操作按钮（注释状态） -->
        <!-- <div class="edit">
          <el-icon @click="toDetail(item)" title="详情"><edit /></el-icon>
          <el-icon title="排行" @click="toRank(item)"><histogram /></el-icon>
          <el-icon title="提交列表" @click="toSubList(item)"><list /></el-icon>
        </div> -->
      </div>
    </div>
    
    <!-- 分页组件 -->
    <div class="pagi">
      <el-pagination
        v-model:currentPage="currentPage" 
        v-model:page-size="pageSize"       
        :page-sizes="[10, 20, 50, 100]"   
        layout="total,sizes, prev, pager, next"
        :total="total"                     
        @size-change="handleSizeChange"    
        @current-change="handleCurrentChange" 
      />
    </div>
  </div>
</template>

<script lang="ts" setup>
// 导入Vue响应式API
import { reactive, ref } from '@vue/reactivity'

// 导入Element Plus图标组件
import { Edit, Histogram, List } from '@element-plus/icons-vue'

// 导入API接口
import api from '../../api/api.js'

// 导入路由功能
import { useRouter } from 'vue-router'
const router = useRouter()  // 创建路由实例

// 响应式数据定义
const rankList = ref([])    // 排行榜列表数据
const sortList = ref([])    // 排序列表数据（未使用）
const pageSize = ref(10)    // 每页显示条数
const currentPage = ref(1)  // 当前页码
const total = ref(0)        // 数据总数

// 分页大小改变处理函数
const handleSizeChange = (val: number) => {
  getRankList()  // 重新获取排行榜数据
}

// 当前页码改变处理函数
const handleCurrentChange = (val: number) => {
  getRankList()  // 重新获取排行榜数据
}

// 获取排行榜数据函数
const getRankList = () => {
  api.getRankList({
    size: pageSize.value,   // 每页大小
    page: currentPage.value // 当前页码
  }).then(res => {
    if (res && res.data) {
      rankList.value = res.data.data.list  // 更新排行榜数据
      total.value = res.data.data.count    // 更新数据总数
    }
    console.log(res)  // 调试输出
  })
}

// 初始化获取排行榜数据
getRankList()

// 查看详情函数（未使用）
const toDetail = (item: any) => {
  router.push({
    path: '/questionDetail',
    query: item  // 传递当前项数据
  })
}

// 查看排行函数（未使用）
const toRank = (item: any) => {
  router.push({
    path: '/topList',
    query: { identity: item.identity }  // 传递身份标识
  })
}

// 查看提交列表函数（未使用）
const toSubList = (item: any) => {
  router.push({
    path: '/submitList',
    query: { identity: item.identity }  // 传递身份标识
  })
}
</script>

<style scoped lang="scss">
/* 排行榜整体样式 */
.top-list {
  height: 100%;                /* 填充父容器高度 */
  display: flex;                /* 弹性布局 */
  justify-content: space-between; /* 两端对齐 */
  flex-direction: column;       /* 垂直方向排列 */
}

/* 列表区域样式 */
.list {
  flex: 1;     /* 占据剩余空间 */
  overflow: auto; /* 内容溢出时滚动 */
}

/* 分页组件样式 */
.pagi {
  text-align: center;       /* 居中对齐 */
  display: flex;            /* 弹性布局 */
  justify-content: center;  /* 水平居中 */
  padding: 10px 0;          /* 上下内边距 */
  border-top: 1px solid #eee; /* 顶部边框 */
}

/* 单个排行榜项样式 */
.item {
  padding: 20px;              /* 内边距 */
  border-bottom: 1px solid #eee; /* 底部边框 */
  display: flex;              /* 弹性布局 */
  justify-content: space-between; /* 两端对齐 */
  align-items: center;        /* 垂直居中 */
 
  /* 标题样式 */
  .title {
    margin-bottom: 20px;      /* 底部外边距 */
    display: flex;            /* 弹性布局 */
    align-items: center;      /* 垂直居中 */
    font-size: 18px;          /* 字体大小 */
    
    /* 标签样式 */
    span {
      background-color: #d3ebff; /* 背景色 */
      font-size: 12px;         /* 较小字体 */
      margin: 0 10px;           /* 左右外边距 */
      padding: 4px;             /* 内边距 */
      border-radius: 5px;       /* 圆角 */
      border: 1px solid #62a9ff; /* 边框 */
      color: #62a9ff;           /* 文字颜色 */
    }
  }
  
  /* 信息区域样式 */
  .msg {
    font-size: 14px;           /* 较小字体 */
    color: #999;               /* 灰色文字 */
    
    /* 信息项样式 */
    span {
      margin-right: 40px;      /* 右侧外边距 */
    }
  }
}

/* 排序列表样式（未使用） */
.sort-list {
  display: flex;              /* 弹性布局 */
  background-color: #d6fff2;  /* 浅绿色背景 */
  padding: 20px;              /* 内边距 */
  
  /* 排序项样式 */
  span {
    color: #13b6f9;           /* 蓝色文字 */
    margin-right: 20px;       /* 右侧外边距 */
    cursor: pointer;          /* 鼠标手型 */
    
    /* 激活状态样式 */
    &.act {
      color: black;           /* 黑色文字 */
    }
  }
}
</style>