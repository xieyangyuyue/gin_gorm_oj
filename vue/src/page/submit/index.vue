<template>
  <!-- 提交列表容器 -->
  <div class="submit-list">
    <!-- 列表区域（带加载状态） -->
    <div class="list" v-loading="loading">
      <!-- 标题行 -->
      <div class="msg title">
        <span>问题</span><span>用户</span><span>提交时间</span>
        <!-- 状态筛选器 -->
        <span style="display:flex;white-space:nowrap;align-items:center">
          状态：
          <!-- 状态选择下拉框 -->
          <el-select v-model="mystatus" clearable @change="getSubmitList">
            <!-- 状态选项 -->
            <el-option :value="i" v-for="(mi,i) in status" :key="mi" :label="mi">{{mi}}</el-option>
          </el-select>
        </span>
      </div>
      
      <!-- 提交列表项 -->
      <div class="msg" v-for="item in submitList" :key="item.id">
        <!-- 问题标题（可点击） -->
        <span @click="toProblem(item.problem_basic)" class="name">{{item.problem_basic.title}}</span>
        <!-- 提交用户 -->
        <span v-if="item.user_basic">{{item.user_basic.name}}</span>
        <!-- 提交时间 -->
        <span v-if="item.user_basic">{{item.created_at}}</span>
        <!-- 提交状态 -->
        <span>{{status[item.status]}}</span>
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

// 导入Element Plus图标组件（未使用）
import { Edit, Histogram, List } from '@element-plus/icons-vue'

// 导入API接口
import api from '../../api/api.js'

// 导入路由功能
import { useRouter, useRoute } from 'vue-router'
const router = useRouter()  // 创建路由实例
const route = useRoute()    // 获取当前路由信息

// 响应式数据定义
const loading = ref(false)  // 加载状态
const mystatus = ref('')    // 当前筛选的状态
const submitList = ref([])  // 提交列表数据
// 状态映射数组
const status = ref(['未知','答案正确','答案错误','超时','运行超内存','编译错误'])

// 分页相关数据
const pageSize = ref(10)    // 每页显示条数
const currentPage = ref(1)  // 当前页码
const total = ref(0)        // 数据总数

// 分页大小改变处理函数
const handleSizeChange = (val: number) => {
  getSubmitList()  // 重新获取提交列表
}

// 当前页码改变处理函数
const handleCurrentChange = (val: number) => {
  getSubmitList()  // 重新获取提交列表
}

// 获取提交列表函数
const getSubmitList = () => {
  loading.value = true  // 显示加载状态
  api.getSubmitList({
    problem_identity: route.query.identity, // 从路由参数获取问题标识
    size: pageSize.value,   // 每页大小
    page: currentPage.value, // 当前页码
    status: mystatus.value   // 筛选状态
  }).then(res => {
    loading.value = false  // 隐藏加载状态
    if (res && res.data) {
      submitList.value = res.data.data.list  // 更新提交列表
      total.value = res.data.data.count      // 更新数据总数
    }
    console.log(res)  // 调试输出
  })
}

// 初始化获取提交列表
getSubmitList()

// 跳转到问题详情页
const toProblem = (detail: any) => {
  router.push({
    path: '/questionDetail',
    query: detail  // 传递问题详情数据
  })
}
</script>

<style scoped lang="scss">
/* 提交列表整体样式 */
.submit-list {
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

/* 列表项样式（未使用） */
.item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

/* 消息行样式 */
.msg {
  font-size: 14px;         /* 较小字体 */
  padding: 20px;            /* 内边距 */
  display: flex;            /* 弹性布局 */
  color: #999;              /* 灰色文字 */
  border-bottom: 1px solid #eee; /* 底部边框 */
  
  /* 列宽设置 */
  span:nth-child(1) {
    width: 50%; /* 问题标题占50% */
  }
  span:nth-child(2) {
    width: 10%; /* 用户占10% */
  }
  span:nth-child(3) {
    width: 20%; /* 提交时间占20% */
  }
  span:nth-child(4) {
    width: 20%; /* 状态占20% */
  }
  
  /* 标题行特殊样式 */
  &.title {
    border-bottom: 1px solid #0087bf; /* 深蓝色底部边框 */
    border-top: 1px solid #0087bf;    /* 深蓝色顶部边框 */
    align-items: center;             /* 垂直居中 */
  }
  
  /* 问题名称样式 */
  .name {
    color: skyblue;   /* 天蓝色文字 */
    cursor: pointer;  /* 鼠标手型 */
  }
}
</style>