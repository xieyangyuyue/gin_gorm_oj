<script setup lang="ts">
// 导入Vue组合式API
import { computed, ref } from 'vue'
// 导入Vuex和Vue Router
import {useStore} from 'vuex'
// 导入登录页面组件
import LoginPage from '../page/login.vue'
// 导入Element Plus消息组件
import {ElMessage} from 'element-plus'
// 导入Vue Router
import {useRouter} from 'vue-router'
// 导入Element Plus图标
import {
  Fold,Expand,Setting
} from '@element-plus/icons-vue'

// 头像URL常量
const  circleUrl='https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'

// 获取Vuex store实例
const store=useStore()
// 获取Vue Router实例
const router=useRouter()

// 计算折叠状态（从Vuex获取）
const collapse=computed(()=>store.state.collapse)
// 计算登录状态
const isLogin=computed(()=>store.state.isLogin)
// 计算管理员状态
const is_admin=computed(()=>store.state.is_admin)
// 计算用户名
const username=computed(()=>store.state.username)

// 切换菜单折叠状态
function changeMenu(){
  store.commit('changeCollapse',!collapse.value)
}

// 处理下拉菜单命令
const handleCommand = (command: string | number | object) => {
  if(command=='a'){ // 退出登录
    localStorage.clear()
    store.commit('logout')
  }else if(command=='b'){ // 分类管理
    router.push('/questionManage')
  }else if(command=='c'){ // 问题管理（开发中）
    ElMessage('正在开发中')
  }
}

// 控制登录对话框显示状态
const showLogin=ref(false)

// 登录成功回调
const loginSucc=()=>{
   showLogin.value=false
}
</script>

<template>
  <!-- 用户信息区域 -->
  <div>
      <!-- 已登录状态显示用户信息 -->
      <el-dropdown @command="handleCommand" v-if="isLogin">
        <!-- 用户信息展示 -->
        <div class="el-dropdown-link">
          <!-- 用户头像 -->
          <el-avatar :size="25" :src="circleUrl" />
          <!-- 用户名 -->
          <span>您好！{{username}}</span>
          <!-- 设置图标 -->
          <el-icon><setting /></el-icon>
        </div>
        <!-- 下拉菜单 -->
        <template #dropdown>
          <el-dropdown-menu>
            <!-- 管理员入口 -->
            <el-dropdown-item :icon="Plus" command="b" v-if="is_admin">进入管理</el-dropdown-item>
            <!-- 退出登录 -->
            <el-dropdown-item :icon="Plus" command="a">退出登录</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
      
      <!-- 未登录状态显示登录按钮 -->
      <span class="login" @click="showLogin=true" v-else>登录</span>
    
    <!-- 登录对话框 -->
    <el-dialog
      v-model="showLogin"
      title="用户登录/注册"
      width="500px"
      :before-close="handleClose"
    >
      <!-- 登录组件 -->
      <LoginPage @loginSucc="loginSucc" /> 
    </el-dialog>
  </div>
</template>

<style scoped lang="scss">
/* 折叠图标样式 */
.fold{
  padding-left: 10px;
}

/* 下拉菜单链接样式 */
.el-dropdown-link{
   display: flex;
   justify-content: space-between;
   align-items: center;
    padding-right: 20px;
    /* 文字间距 */
    span{
      padding: 0 10px;
    }
 }

/* 登录按钮样式 */
.login{
   color: white;
   padding: 0 20px;
   cursor: pointer;
 }
</style>