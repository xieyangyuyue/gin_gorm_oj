<template>
  <!-- 问题添加/编辑表单 -->
  <el-form
      ref="ruleFormRef"                     
      :model="ruleForm"                    
      :rules="rules"                       
      label-width="120px"                  
      class="demo-ruleForm"                
      :size="formSize"                     
  >
    <!-- 问题名称输入项 -->
    <el-form-item label="问题名称" prop="title">
      <el-input v-model="ruleForm.title"/>  <!-- 双向绑定问题标题 -->
    </el-form-item>
    
    <!-- 问题内容输入项 -->
    <el-form-item label="内容" prop="content ">
      <el-input v-model="ruleForm.content" type="textarea"/> <!-- 多行文本输入 -->
    </el-form-item>
    
    <!-- 分类选择项 -->
    <el-form-item label="分类" prop="problem_categories">
      <el-select
          v-model="ruleForm.problem_categories" 
          placeholder="请选择"
          multiple                            
          style="width: 100%"                 
      >
        <!-- 遍历分类选项 -->
        <el-option
            :label="mi.name"                   
            :value="mi.id"                     
            v-for="mi in sortList"             
            :key="mi.id"                       
        />
      </el-select>
    </el-form-item>
    
    <!-- 最大运行时间输入项 -->
    <el-form-item label="最大运行时间" prop="max_runtime">
      <el-input v-model.number="ruleForm.max_runtime"/> <!-- 数字类型绑定 -->
    </el-form-item>
    
    <!-- 最大内存占用输入项 -->
    <el-form-item label="最大占用内存" prop="max_mem">
      <el-input v-model.number="ruleForm.max_mem"/> <!-- 数字类型绑定 -->
    </el-form-item>

    <!-- 测试案例输入区域 -->
    <el-form-item label="测试案例" prop="test_cases">
      <!-- 遍历测试案例数组 -->
      <el-form-item
          label-width="0"                     
          v-for="(mi, i) in ruleForm.test_cases" 
          :key="i"                        
      >
        <!-- 测试案例行布局 -->
        <el-row :gutter="20">                  <!-- 栅格布局 -->
          <!-- 输入区域 -->
          <el-col :span="4" style="text-align: right">输入：</el-col>
          <el-col :span="6">
            <el-input type="textarea" v-model="mi.input" rows="6" cols="40"/> <!-- 输入框 -->
          </el-col>
          
          <!-- 输出区域 -->
          <el-col :span="4" style="text-align: right">输出：</el-col>
          <el-col :span="6">
            <el-input type="textarea" v-model="mi.output" rows="6" cols="40"/> <!-- 输出框 -->
          </el-col>
          
          <!-- 操作按钮 -->
          <el-col :span="4">
            <!-- 添加案例按钮 -->
            <el-icon @click="addCase">
              <circle-plus/>
            </el-icon>
            <!-- 删除案例按钮 -->
            <el-icon @click="removeCase(i)">
              <remove/>
            </el-icon>
          </el-col>
        </el-row>
      </el-form-item>
    </el-form-item>
    
    <!-- 表单操作按钮 -->
    <el-form-item>
      <el-button type="primary" @click="submitForm(ruleFormRef)">创建</el-button>
      <el-button @click="closeBox">取消</el-button>
    </el-form-item>
  </el-form>
</template>

<script lang="ts" setup>
// 导入Vue相关功能
import {reactive, ref} from "vue";

// 导入Element Plus表单相关类型
import type {FormInstance, FormRules} from "element-plus";

// 导入Element Plus消息组件
import {ElMessage} from "element-plus";

// 导入Element Plus图标组件
import {CirclePlus, Remove} from "@element-plus/icons-vue";

// 导入API接口
import api from "../../api/api.js";

// 表单尺寸控制
const formSize = ref("default");

// 表单引用
const ruleFormRef = ref<FormInstance>();

// 表单数据对象（响应式）
const ruleForm = ref({
  title: "",            // 问题标题
  content: "",          // 问题内容
  max_runtime: 0,       // 最大运行时间
  max_mem: 0,           // 最大内存占用
  test_cases: [{input: "", output: ""}], // 测试案例数组（至少一个）
  problem_categories: [], // 分类ID数组
});

// 定义组件事件
const emits = defineEmits(["cancel", "submit"]);

// 定义组件属性
const props = defineProps(["sortList", "selectQues"]);

// 如果传入选择的问题（编辑模式），初始化表单数据
if (props.selectQues) {
  let se = props.selectQues;  // 当前选择的问题
  
  // 处理分类数据
  if (se.problem_categories instanceof Array) {
    se.problem_categories = se.problem_categories.map((it: any) => {
      return it.category_id;  // 提取分类ID
    });
  }
  
  // 处理测试案例数据
  if (!Array.isArray(se.test_cases) || se.test_cases.length == 0) {
    se.test_cases = [{input: "", output: ""}]; // 确保至少一个测试案例
  }
  
  // 更新表单数据
  ruleForm.value = se;
}

// 关闭表单函数
const closeBox = () => {
  emits("cancel");  // 触发取消事件
};

// 表单验证规则
const rules = reactive<FormRules>({
  title: [{required: true, message: "请输入", trigger: "blur"}], // 标题必填
  content: [{required: true, message: "请输入", trigger: "blur"}], // 内容必填
  max_runtime: [
    {required: true, message: "请输入", trigger: "blur"}, // 必填
    {type: "number", message: "请输入数字", trigger: "blur"}, // 必须为数字
  ],
  max_mem: [
    {required: true, message: "请输入", trigger: "blur"}, // 必填
    {type: "number", message: "请输入数字", trigger: "blur"}, // 必须为数字
  ],
  test_cases: [{required: true, message: "请输入", trigger: "blur"}], // 至少一个测试案例
  problem_categories: [
    {
      type: "array", // 数组类型
      required: true, // 必填
      message: "请至少选择一个分类", // 错误消息
      trigger: "change", // 触发时机
    },
  ],
});

// 添加测试案例函数
const addCase = () => {
  ruleForm.value.test_cases.push({
    input: "",  // 空输入
    output: "", // 空输出
  });
};

// 删除测试案例函数
const removeCase = (i: number) => {
  // 确保至少保留一个测试案例
  if (ruleForm.value.test_cases.length > 1) {
    ruleForm.value.test_cases.splice(i, 1); // 删除指定索引的案例
  }
};

// 提交表单函数
const submitForm = async (formEl: FormInstance | undefined) => {
  if (!formEl) return; // 表单引用检查
  
  // 表单验证
  await formEl.validate((valid, fields) => {
    if (valid) { // 验证通过
      // 编辑模式
      if (ruleForm.value.identity) {
        api.editProblem(ruleForm.value).then((res) => {
          if (res.data.code == 200) {
            ElMessage.success("成功"); // 成功提示
            emits("submit"); // 触发提交事件
          } else {
            ElMessage(res.data.msg); // 错误提示
          }
        });
      } else { // 新增模式
        api.addProblem(ruleForm.value).then((res: any) => {
          if (res.data.code == 200) {
            ElMessage.success("成功"); // 成功提示
            emits("submit"); // 触发提交事件
          } else {
            ElMessage(res.data.msg); // 错误提示
          }
        });
      }
    } else { // 验证失败
      console.log("error submit!", fields); // 调试输出
    }
  });
};

// 重置表单函数（未使用）
const resetForm = (formEl: FormInstance | undefined) => {
  if (!formEl) return;
  formEl.resetFields(); // 重置表单字段
};
</script>