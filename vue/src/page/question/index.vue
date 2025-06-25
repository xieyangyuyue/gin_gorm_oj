<template>
  <!-- 问题列表容器 -->
  <div class="ques-list">
    <!-- 分类筛选区域 -->
    <div class="sort-list">
      <!-- 问题分类标签 -->
      <div>
        问题分类：
        <!-- 全部分类标签 -->
        <span @click="getProblem(null)" :class="!actSort ? 'act' : ''"
          >全部</span
        >
        <!-- 遍历分类列表 -->
        <span
          v-for="item in sortList"
          :key="item.id"
          @click="getProblem(item.identity)"
          :class="actSort == item.identity ? 'act' : ''"
          >{{ item.name }}</span
        >
      </div>
      <!-- 搜索输入框 -->
      <el-input
        style="width:200px;margin-right:10px"
        v-model="keyword"            
        clearable                 
        @change="getProblem(actSort)" 
        @clear="getProblem(null)"    
        size="large"                
        placeholder="请搜索"         
        :suffix-icon="Search"        
      />
    </div>
    
    <!-- 问题列表区域（带加载状态） -->
    <div class="list" v-loading="loading">
      <!-- 遍历问题列表 -->
      <div class="item" v-for="item in quesList" :key="item.id">
        <div>
          <!-- 问题标题和分类 -->
          <div class="title">
            <!-- 可点击的问题标题 -->
            <b @click="toDetail(item)" >{{ item.title }}</b> 
            <!-- 问题分类标签 -->
            <div class="sort">
              <span v-for="mi in item.problem_categories" :key="mi.id">
                <!-- 显示分类名称 -->
                <b v-if="mi.category_basic">{{ mi.category_basic.name }}</b> 
              </span>
            </div>
          </div>
          <!-- 问题元信息 -->
          <div class="msg">
            <span>创建时间：{{ item.created_at }}</span> <!-- 创建时间 -->
            <span>提交人数：{{ item.submit_num }}</span> <!-- 提交人数 -->
            <span>通过人数：{{ item.pass_num }}</span>   <!-- 通过人数 -->
          </div>
        </div>
        <!-- 操作按钮区域 -->
        <div class="edit">
          <!-- 详情按钮 -->
          <el-icon @click="toDetail(item)" title="详情"><edit /></el-icon>
          <!-- 排行按钮（注释状态） -->
          <!-- <el-icon title="排行" @click="toRank(item)"><histogram /></el-icon> -->
          <!-- 提交列表按钮 -->
          <el-icon title="提交列表" @click="toSubList(item)"><list /></el-icon>
        </div>
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
import { reactive, ref } from "@vue/reactivity";

// 导入Element Plus图标组件
import { Edit, Histogram, List } from "@element-plus/icons-vue";

// 导入API接口
import api from "../../api/api.js";

// 导入路由功能
import { useRouter } from "vue-router";
const router = useRouter();  // 创建路由实例

// 响应式数据定义
const loading = ref(false)   // 加载状态
const quesList = ref([]);    // 问题列表数据
const sortList = ref([]);    // 分类列表数据
const pageSize = ref(10)     // 每页显示条数
const currentPage = ref(1)   // 当前页码
const total = ref(0)         // 数据总数
const actSort = ref<null | number>(null);  // 当前激活的分类ID
const keyword = ref('')      // 搜索关键词

// 分页大小改变处理函数
const handleSizeChange = (val: number) => {
  getProblem(actSort.value)  // 重新获取问题列表
}

// 当前页码改变处理函数
const handleCurrentChange = (val: number) => {
  getProblem(actSort.value)  // 重新获取问题列表
}

// 获取问题列表函数
const getProblem = (sortId: number | null) => {
  actSort.value = sortId;    // 更新激活分类
  loading.value = true       // 显示加载状态
  api
    .getProblemList({
      category_identity: sortId,  // 分类标识
      size: pageSize.value,       // 每页大小
      page: currentPage.value,    // 当前页码
      keyword: keyword.value      // 搜索关键词
    })
    .then((res) => {
      loading.value = false  // 隐藏加载状态
      if (res && res.data) {
        quesList.value = res.data.data.list;  // 更新问题列表
        total.value = res.data.data.count     // 更新数据总数
      }
      console.log(res);  // 调试输出
    });
};

// 初始化获取问题列表（全部问题）
getProblem(null);

// 获取分类列表
api.getSortList({}).then((res) => {
  if (res && res.data) {
    sortList.value = res.data.data.list;  // 更新分类列表
  }
});

// 跳转到问题详情页
const toDetail = (item: any) => {
  router.push({
    path: "/questionDetail",
    query: item,  // 传递问题数据
  });
};

// 跳转到排行页（未使用）
const toRank = (item: any) => {
  router.push({
    path: "/topList",
    query: { identity: item.identity },  // 传递问题标识
  });
};

// 跳转到提交列表页
const toSubList = (item: any) => {
  router.push({
    path: "/submitList",
    query: { identity: item.identity },  // 传递问题标识
  });
};
</script>

<style scoped lang="scss">
/* 问题列表整体样式 */
.ques-list {
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

/* 单个问题项样式 */
.item {
  padding: 20px;              /* 内边距 */
  border-bottom: 1px solid #eee; /* 底部边框 */
  display: flex;              /* 弹性布局 */
  justify-content: space-between; /* 两端对齐 */
  align-items: center;        /* 垂直居中 */
  
  /* 操作按钮区域 */
  .edit {
    width: 80px;             /* 固定宽度 */
    display: flex;           /* 弹性布局 */
    justify-content: space-between; /* 按钮均匀分布 */
    color: #13b6f9;          /* 蓝色图标 */
    cursor: pointer;         /* 鼠标手型 */
  }
  
  /* 标题区域 */
  .title {
    margin-bottom: 20px;      /* 底部外边距 */
    display: flex;            /* 弹性布局 */
    align-items: center;      /* 垂直居中 */
    font-size: 18px;          /* 较大字体 */
    
    /* 问题标题样式 */
    b {
      cursor: pointer;        /* 鼠标手型 */
    }
    
    /* 分类标签样式 */
    span {
      background-color: #d3ebff; /* 浅蓝色背景 */
      font-size: 12px;         /* 较小字体 */
      margin: 0 10px;           /* 左右外边距 */
      padding: 4px;             /* 内边距 */
      border-radius: 5px;       /* 圆角 */
      border: 1px solid #62a9ff; /* 边框 */
      color: #62a9ff;           /* 文字颜色 */
    }
  }
  
  /* 元信息区域 */
  .msg {
    font-size: 14px;           /* 较小字体 */
    color: #999;               /* 灰色文字 */
    
    /* 信息项样式 */
    span {
      margin-right: 20px;      /* 右侧外边距 */
    }
  }
}

/* 分类筛选区域样式 */
.sort-list {
  display: flex;              /* 弹性布局 */
  align-items: center;        /* 垂直居中 */
  justify-content: space-between; /* 两端对齐 */
  border-bottom: 1px solid #0087bf; /* 深蓝色底部边框 */
  border-top: 1px solid #0087bf;    /* 深蓝色顶部边框 */
  padding: 20px;              /* 内边距 */
  
  /* 分类标签样式 */
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