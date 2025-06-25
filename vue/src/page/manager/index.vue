<template>
  <!-- 问题管理页面容器 -->
  <div class="ques-list">
    <!-- 分类管理对话框 -->
    <el-dialog
      v-model="sortDialog"                    
      :title="selectSort.id ? '编辑分类' : '新增分类'" 
      width="30%"                           
      :before-close="closeSort"             
    >
      <!-- 分类名称输入框 -->
      <el-input
        v-model="selectSort.name"          
        placeholder="请输入分类名称"      
      ></el-input>
      <!-- 对话框底部按钮 -->
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="closeSort">取消</el-button>
          <el-button type="primary" @click="subSort">确定</el-button>
        </span>
      </template>
    </el-dialog>
    
    <!-- 问题管理对话框 -->
    <el-dialog
      v-model="quesDialog"                    
      :title="selectQues.id ? '编辑问题' : '新增问题'" 
      width="70%"                             
      :close-on-click-modal="false"           
      :before-close="closeQues"               
    >
      <!-- 添加/编辑问题组件 -->
      <addQues 
        @cancel="closeQues"                    
        :selectQues="selectQues"            
        @submit="subQues"               
        :sortList="sortList"               
        v-if="quesDialog"               
      ></addQues> 
    </el-dialog>
    
    <!-- 主体内容区域 -->
    <div class="flex-box">
      <!-- 分类列表区域 -->
      <div class="sort-list">
        <!-- 分类标题 -->
        <h3>
          分类：<el-icon @click="sortDialog = true"><circle-plus /></el-icon> <!-- 添加分类按钮 -->
        </h3>
        
        <!-- 全部分类项 -->
        <div @click="getProblem(null)" :class="!actSort ? 'act' : ''">全部</div>
        
        <!-- 遍历分类列表 -->
        <div
          v-for="item in sortList"
          :key="item.id"
          @click="getProblem(item.identity)"
          :class="actSort == item.identity ? 'act' : ''"
        >
          <!-- 分类名称 -->
          <span>{{ item.name }}</span>
          
          <!-- 分类操作按钮 -->
          <div class="edit">
            <el-icon title="删除" @click.stop="delSort(item)"><delete /></el-icon>
            <el-icon title="编辑" @click.stop="editSort(item)"><edit /></el-icon>
          </div>
        </div>
      </div>
      
      <!-- 问题列表区域 -->
      <div class="list" v-loading="loading"> <!-- 加载状态 -->
        <!-- 列表标题和操作区 -->
        <h3>
          <!-- 搜索框 -->
          <el-input
            style="width: 200px; margin-right: 10px"
            v-model="keyword"                 
            clearable                         
            @change="getProblem(actSort)"     
            @clear="getProblem(null)"         
            placeholder="请搜索"              
            :suffix-icon="Search"             
          />
          <!-- 新增问题按钮 -->
          <el-button type="primary" @click="quesDialog=true">新增问题</el-button>
        </h3>
        
        <!-- 滚动容器 -->
        <div class="scroll">
          <!-- 遍历问题列表 -->
          <div class="item" v-for="item in quesList" :key="item.id">
            <!-- 问题信息 -->
            <div>
              <!-- 问题标题和分类 -->
              <div class="title">
                <b>{{ item.title }}</b> <!-- 问题标题 -->
                <div class="sort">
                  <!-- 遍历问题所属分类 -->
                  <span v-for="mi in item.problem_categories" :key="mi.id" >
                    <b v-if="mi.category_basic">{{ mi.category_basic.name }}</b> <!-- 分类名称 -->
                  </span>
                </div>
              </div>
              <!-- 问题元数据 -->
              <div class="msg">
                <span>创建时间：{{ item.created_at }}</span> <!-- 创建时间 -->
                <span>提交人数：{{ item.submit_num }}</span> <!-- 提交人数 -->
                <span>通过人数：{{ item.pass_num }}</span>   <!-- 通过人数 -->
              </div>
            </div>
            <!-- 问题操作按钮 -->
            <div class="edit">
              <el-icon title="编辑" @click="editQues(item)"><edit /></el-icon> <!-- 编辑按钮 -->
            </div>
          </div>
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
import { Edit, Delete, List, CirclePlus } from "@element-plus/icons-vue";

// 导入Element Plus消息组件
import { ElMessageBox, ElMessage } from "element-plus";

// 导入API接口
import api from "../../api/api.js";

// 导入路由功能
import { useRouter } from "vue-router";

// 导入添加/编辑问题组件
import addQues from './add.vue'

// 创建路由实例
const router = useRouter();

// 响应式数据定义
const loading = ref(false);          // 加载状态
const quesList = ref([]);            // 问题列表数据
const sortList = ref([]);            // 分类列表数据
const sortDialog = ref(false);       // 分类对话框显示状态
const selectSort = ref({ name: "",id:'' }); // 当前选中的分类

const quesDialog=ref(false)          // 问题对话框显示状态
const selectQues=ref({})             // 当前选中的问题

const pageSize = ref(10);            // 每页显示条数
const currentPage = ref(1);          // 当前页码
const total = ref(0);                // 数据总数
const actSort = ref<null | number>(null); // 当前激活的分类ID
const keyword = ref("");             // 搜索关键词

// 关闭问题对话框函数
const closeQues=()=>{
  quesDialog.value=false             // 关闭对话框
  selectQues.value={}                // 清空选中问题
}

// 问题提交成功回调
const subQues=()=>{
  closeQues()                        // 关闭对话框
  getProblem(actSort.value)          // 刷新问题列表
}

// 编辑问题函数
const editQues=(item:any)=>{
  // 获取问题详情
  api.getProblemDetail({
      identity:item.identity
  }).then(res=>{
        if(res.data.data){
            selectQues.value=res.data.data // 设置选中问题
            quesDialog.value=true          // 打开对话框
        }
  })
// selectQues.value=item
}

// 分页大小改变处理函数
const handleSizeChange = (val: number) => {
  getProblem(actSort.value);         // 重新获取问题列表
};

// 编辑分类函数
const editSort = (item: any) => {
  selectSort.value = {...item};       // 设置选中分类
  sortDialog.value = true;            // 打开对话框
};

// 当前页码改变处理函数
const handleCurrentChange = (val: number) => {
  getProblem(actSort.value);         // 重新获取问题列表
};

// 关闭分类对话框函数
const closeSort = () => {
  sortDialog.value = false;           // 关闭对话框
  selectSort.value = { name: "",id:'' }; // 清空选中分类
};

// 删除分类函数
const delSort = (item: any) => {
  // 确认对话框
  ElMessageBox.confirm("确定要删除该分类吗?", "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  })
    .then(() => {
      // 确认删除
      api.delSort({ identity: item.identity }).then((res) => {
        if (res.data.code == 200) {
          ElMessage.success(res.data.msg); // 成功提示
          getSortList();                   // 刷新分类列表
        } else {
          ElMessage.warning(res.data.msg); // 警告提示
        }
      });
    })
    .catch(() => {
      // 取消操作
    });
};

// 获取问题列表函数
const getProblem = (sortId: number | null) => {
  actSort.value = sortId;             // 更新激活分类
  loading.value = true;               // 显示加载状态
  api
    .getProblemList({
      category_identity: sortId,      // 分类标识
      size: pageSize.value,           // 每页大小
      page: currentPage.value,        // 当前页码
      keyword: keyword.value,         // 搜索关键词
    })
    .then((res) => {
      loading.value = false;          // 隐藏加载状态
      if (res && res.data) {
        quesList.value = res.data.data.list; // 更新问题列表
        total.value = res.data.data.count;   // 更新数据总数
      }
      console.log(res);               // 调试输出
    });
};

// 初始化获取问题列表（全部问题）
getProblem(null);

// 获取分类列表函数
const getSortList = () => {
  api.getSortList({}).then((res) => {
    if (res && res.data) {
      sortList.value = res.data.data.list; // 更新分类列表
    }
  });
};

// 初始化获取分类列表
getSortList();

// 提交分类函数（新增/编辑）
const subSort = () => {
  if (selectSort.value.name) {        // 验证分类名称
    if (selectSort.value.id) {
      // 编辑分类
      api.editSort(selectSort.value).then((res) => {
        if (res.data.code == 200) {
          ElMessage.success(res.data.msg); // 成功提示
          closeSort();                     // 关闭对话框
          getSortList()                    // 刷新分类列表
        } else {
          ElMessage.warning(res.data.msg); // 警告提示
        }
      });
    } else {
      // 新增分类
      api.addSort({ name: selectSort.value.name }).then((res) => {
        if (res.data.code == 200) {
          ElMessage.success(res.data.msg); // 成功提示
          closeSort();                     // 关闭对话框
          getSortList()                    // 刷新分类列表
        } else {
          ElMessage.warning(res.data.msg); // 警告提示
        }
      });
    }
  } else {
    ElMessage.warning("请输入分类名称"); // 验证失败提示
  }
};

// 跳转到问题详情页（未使用）
const toDetail = (item: any) => {
  router.push({
    path: "/questionDetail",
    query: item,
  });
};
</script>

<style scoped lang="scss">
/* 滚动容器样式 */
.scroll {
  height: calc(100% - 60px); /* 计算高度 */
  overflow: auto;            /* 允许滚动 */
}

/* 问题列表整体样式 */
.ques-list {
  height: 100%;              /* 100%高度 */
  display: flex;             /* 弹性布局 */
  justify-content: space-between; /* 两端对齐 */
  flex-direction: column;    /* 垂直排列 */
}

/* 弹性盒子布局 */
.flex-box {
  display: flex;             /* 弹性布局 */
  height: calc(100% - 60px); /* 计算高度 */
}

/* 列表区域样式 */
.list {
  flex: 1;                   /* 占据剩余空间 */
  display: flex;              /* 弹性布局 */
  flex-direction: column;     /* 垂直排列 */
  height: 100%;               /* 100%高度 */
  
  /* 标题样式 */
  h3 {
    display: flex;            /* 弹性布局 */
    justify-content: space-between; /* 两端对齐 */
    align-items: center;      /* 垂直居中 */
    padding: 10px;            /* 内边距 */
  }
}

/* 分页组件样式 */
.pagi {
  text-align: center;         /* 居中对齐 */
  display: flex;              /* 弹性布局 */
  justify-content: center;    /* 水平居中 */
  padding: 10px 0;            /* 上下内边距 */
  border-top: 1px solid #eee; /* 顶部边框 */
}

/* 问题项样式 */
.item {
  padding: 20px;              /* 内边距 */
  border-bottom: 1px solid #eee; /* 底部边框 */
  display: flex;              /* 弹性布局 */
  justify-content: space-between; /* 两端对齐 */
  align-items: center;        /* 垂直居中 */
  
  /* 编辑按钮区域 */
  .edit {
    width: 80px;             /* 固定宽度 */
    display: flex;            /* 弹性布局 */
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
    
    /* 问题标题 */
    b {
      cursor: pointer;        /* 鼠标手型 */
    }
    
    /* 分类标签 */
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
  
  /* 元数据区域 */
  .msg {
    font-size: 14px;           /* 较小字体 */
    color: #999;               /* 灰色文字 */
    
    /* 数据项样式 */
    span {
      margin-right: 20px;      /* 右侧外边距 */
    }
  }
}

/* 分类列表区域样式 */
.sort-list {
  width: 250px;               /* 固定宽度 */
  align-items: center;        /* 垂直居中 */
  justify-content: space-between; /* 两端对齐 */
  border-right: 1px solid #0087bf; /* 右侧边框 */
  
  /* 分类项样式 */
  & > div {
    display: flex;            /* 弹性布局 */
    justify-content: space-between; /* 两端对齐 */
    padding: 10px;            /* 内边距 */
  }
  
  /* 编辑按钮区域 */
  .edit {
    width: 20%;               /* 相对宽度 */
    display: flex;            /* 弹性布局 */
    justify-content: space-between; /* 按钮均匀分布 */
  }
  
  /* 标题样式 */
  h3 {
    border-bottom: 1px solid #0087bf; /* 底部边框 */
    display: flex;            /* 弹性布局 */
    justify-content: space-between; /* 两端对齐 */
    align-items: center;      /* 垂直居中 */
    padding: 14px;            /* 内边距 */
  }
  
  /* 分类项通用样式 */
  div {
    color: #13b6f9;          /* 蓝色文字 */
    cursor: pointer;          /* 鼠标手型 */
    
    /* 激活状态样式 */
    &.act {
      background-color: #c4e9f9; /* 浅蓝色背景 */
    }
  }
}
</style>