<template>
  <!-- 登录页面容器 -->
  <div class="login-page">
    <!-- Element Plus 标签页组件 -->
    <el-tabs
        v-model="activeName"
        type="card"
        class="demo-tabs"
        @tab-click="handleClick"
    >
      <!-- 第一个标签页：登录 -->
      <el-tab-pane label="登录" name="first">
        <!-- 登录表单组件 -->
        <el-form
            ref="ruleFormRef"
            :model="ruleForm"
            :rules="rules"
            label-width="80px"
            class="login-form"
            :size="formSize"
        >
          <!-- 用户名输入项 -->
          <el-form-item label="账户" prop="username">
            <el-input v-model="ruleForm.username"/>
          </el-form-item>

          <!-- 密码输入项 -->
          <el-form-item label="密码" prop="password">
            <el-input v-model="ruleForm.password" show-password/>
          </el-form-item>
          <!-- 登录按钮 -->
          <div style="text-align: center">
            <el-button type="primary" @click="submitForm(ruleFormRef)"
            >登录
            </el-button
            >
          </div>
        </el-form>
      </el-tab-pane>
      
      <!-- 第二个标签页：注册 -->
      <el-tab-pane label="注册" name="second">
        <!-- 注册表单组件 -->
        <el-form
            ref="registFormRef"
            :model="registForm"
            :rules="rules"
            label-width="80px"
            class="login-form"
            :size="formSize"
        >
          <!-- 用户名输入 -->
          <el-form-item label="账户" prop="name">
            <el-input v-model="registForm.name"/>
          </el-form-item>

          <!-- 密码输入 -->
          <el-form-item label="密码" prop="password">
            <el-input v-model="registForm.password" show-password/>
          </el-form-item>
          
          <!-- 邮箱输入 -->
          <el-form-item label="邮箱" prop="mail">
            <el-input v-model="registForm.mail"/>
          </el-form-item>
          
          <!-- 验证码输入（带发送按钮） -->
          <el-form-item label="验证码" prop="code">
            <el-row :gutter="20">
              <el-col :span="12">
                <el-input v-model="registForm.code"/>
              </el-col>
              <el-col :span="12" style="text-align:right;">
                <!-- 倒计时显示按钮 -->
                <el-button disabled v-if="remainTime>0&&remainTime<60">{{ remainTime }}秒</el-button>
                <!-- 正常发送验证码按钮 -->
                <el-button @click="sendCode" v-else type="primary">发送验证码</el-button>
              </el-col>
            </el-row>
          </el-form-item>

          <!-- 注册按钮 -->
          <div style="text-align: center">
            <el-button type="primary" @click="subRegister(registFormRef)"
            >注册
            </el-button
            >
          </div>
        </el-form>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script lang="ts" setup>
// Vue 组合式 API
import { reactive, ref } from "vue";
// 状态管理
import { useStore } from "vuex";
// 路由控制
import { useRouter } from "vue-router";
// Element Plus 类型引入
import type { FormInstance, TabsPaneContext } from "element-plus";
// Element Plus 消息组件
import { ElMessage } from "element-plus";
// API 接口
import api from '../api/api.js';

// 当前激活的标签页（默认第一个）
const activeName = ref("first");

// 标签页切换处理函数
const handleClick = (tab: TabsPaneContext, event: Event) => {
  console.log(tab, event);  // 打印切换信息（调试用）
};

// 定义事件发射器（用于登录成功通知父组件）
const emits = defineEmits(["loginSucc"]);

// 表单尺寸控制
const formSize = ref("default");
// Vuex 状态管理实例
const store = useStore();

// 表单引用（用于验证）
const ruleFormRef = ref<FormInstance>();  // 登录表单引用
const registFormRef = ref<FormInstance>() // 注册表单引用
const remainTime = ref(60)                // 验证码倒计时（秒）

// 登录表单数据
const ruleForm = reactive({
  username: "",
  password: "",
});

// 注册表单数据
const registForm = reactive({
  name: "",
  password: "",
  mail: '',
  code: ''
});

// 表单验证规则
const rules = reactive({
  username: [
    { required: true, message: "请输入用户名", trigger: "blur" },
  ],
  name: [  // 注册用户名验证
    { required: true, message: "请输入用户名", trigger: "blur" },
  ],
  code: [  // 验证码验证
    { required: true, message: "请输入验证码", trigger: "blur" },
  ],
  mail: [  // 邮箱验证
    { required: true, message: "请输入邮箱", trigger: "blur" },
  ],
  password: [  // 密码验证（登录和注册共用）
    { required: true, message: "请输入密码", trigger: "blur" },
  ],
});

// 路由实例
const router = useRouter();
console.log(router);  // 调试路由信息

// 登录表单提交
const submitForm = async (formEl: FormInstance | undefined) => {
  if (!formEl) return;  // 空引用检查
  
  // 表单验证
  await formEl.validate((valid, fields) => {
    if (valid) {  // 验证通过
      // 调用登录API
      api.login(ruleForm).then((res: any) => {
        if (res.data.code == 200) {  // 登录成功
          ElMessage.success('登录成功');
          
          // 存储 token 到 localStorage
          localStorage.setItem("token", res.data.data.token);
          // 更新 Vuex 状态
          store.commit("loginSucc", res.data.data.token);
          store.commit("setUser", {
            username: ruleForm.username,
            is_admin: res.data.data.is_admin
          });
          
          // 存储用户信息到 localStorage
          localStorage.setItem('is_admin', res.data.data.is_admin);
          localStorage.setItem('username', ruleForm.username);

          // 通知父组件登录成功
          emits("loginSucc");
          // 跳转到首页
          router.push('/index');
        } else {  // 登录失败
          ElMessage.error(res.data.msg);
        }
      });
    } else {  // 验证失败
      console.log("error submit!", fields);
    }
  });
};

// 注册表单提交
const subRegister = async (formEl: FormInstance | undefined) => {
  if (!formEl) return;
  
  await formEl.validate((valid, fields) => {
    if (valid) {  // 验证通过
      // 调用注册API
      api.register(registForm).then((res:any) => {
        if (res.data.code == 200) {  // 注册成功
          ElMessage.success('注册成功');
          
          // 存储 token 和用户信息（自动登录）
          localStorage.setItem("token", res.data.data.token);
          store.commit("loginSucc", res.data.data.token);
          store.commit("setUser", {
            username: registForm.name,
            is_admin: res.data.data.is_admin
          });
          localStorage.setItem('username', registForm.name);
          localStorage.setItem('is_admin', res.data.data.is_admin);

          // 通知父组件并跳转
          emits("loginSucc");
          router.push('/index');
        } else {  // 注册失败
          ElMessage.error(res.data.msg);
        }
      });
    } else {  // 验证失败
      console.log("error submit!", fields);
    }
  });
};

// 表单重置函数（未使用）
const resetForm = (formEl: FormInstance | undefined) => {
  if (!formEl) return;
  formEl.resetFields();
};

// 倒计时计时器函数
const startRemain = () => {
  if (remainTime.value > 0) {
    remainTime.value--  // 秒数减1
    setTimeout(startRemain, 1000)  // 1秒后递归调用
  } else {
    remainTime.value = 60  // 重置倒计时
  }
}

// 发送验证码
const sendCode = () => {
  // 邮箱正则验证
  const re = /^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$/;
  if (re.test(registForm.mail)) {  // 邮箱格式正确
    startRemain()  // 启动倒计时
    
    // 调用发送验证码API
    api.sendCode({
      email: registForm.mail
    }).then((res: any) => {
      if (res.data.code == 200) {
        ElMessage.success(res.data.msg)  // 发送成功提示
      } else {
        ElMessage.error(res.data.msg)    // 发送失败提示
      }
    })
  } else {  // 邮箱格式错误
    ElMessage('请输入正确的邮箱')
  }
}
</script>

<style lang="scss" scoped>
/* 登录页面样式 */
.login-page {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;  /* 水平居中 */
  align-items: center;      /* 垂直居中 */
}

/* 登录/注册表单样式 */
.login-form {
  border: 1px solid #eee;     /* 浅灰色边框 */
  padding: 40px 80px 20px 80px;  /* 内边距 */
  border-radius: 10px;        /* 圆角 */
}
</style>