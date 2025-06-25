<template>
  <!-- 代码编辑器容器 -->
  <div class="e-box">
    <!-- 语言选择器 -->
    <div class="select">
      <el-select v-model="language" @change="changeLanguage">
        <el-option value="go">go</el-option> <!-- 目前只支持Go语言 -->
      </el-select>
    </div>
    
    <!-- 代码编辑器挂载点 -->
    <div id="codeEditBox"></div>
    
    <!-- 提交按钮 -->
    <div class="submit">
      <el-button 
        type="primary" 
        @click="submitCode" 
        :loading="loading" 
      >提交</el-button>
    </div>
    
    <!-- 编译结果展示区 -->
    <div class="sub-box">
     编译结果： {{msg}} <!-- 显示编译结果信息 -->
    </div>
  </div>
</template>

<script lang="ts" setup>
// 导入Monaco Editor的各种语言的worker文件
import jsonWorker from 'monaco-editor/esm/vs/language/json/json.worker?worker'
import cssWorker from 'monaco-editor/esm/vs/language/css/css.worker?worker'
import htmlWorker from 'monaco-editor/esm/vs/language/html/html.worker?worker'
import tsWorker from 'monaco-editor/esm/vs/language/typescript/ts.worker?worker'
import EditorWorker from 'monaco-editor/esm/vs/editor/editor.worker?worker'

// 导入Monaco Editor核心库
import * as monaco from 'monaco-editor';

// Vue相关功能导入
import { nextTick, ref, onBeforeUnmount } from 'vue'

// 路由相关
import { useRoute } from 'vue-router'

// API接口
import api from '../../api/api.js'

// Element Plus消息组件
import { ElMessage } from 'element-plus'

// 响应式数据定义
const text = ref('')         // 编辑器中的代码文本
const route = useRoute()     // 获取当前路由信息
const language = ref('go')   // 当前选择的编程语言（默认Go）
const msg = ref()            // 编译结果消息
const loading = ref(false)   // 提交按钮加载状态

// 组件卸载前的清理工作
onBeforeUnmount(() => {
  editor.dispose()  // 销毁编辑器实例
})

// 配置Monaco Editor的工作环境
// @ts-ignore
self.MonacoEnvironment = {
  // 根据语言标签返回对应的worker
  getWorker(_: string, label: string) {
    if (label === 'json') {
      return new jsonWorker()
    }
    if (label === 'css' || label === 'scss' || label === 'less') {
      return new cssWorker()
    }
    if (label === 'html' || label === 'handlebars' || label === 'razor') {
      return new htmlWorker()
    }
    if (['typescript', 'javascript'].includes(label)) {
      return new tsWorker()
    }
    return new EditorWorker()  // 默认worker
  },
}

// 编辑器实例变量
let editor: monaco.editor.IStandaloneCodeEditor;

// 初始化编辑器
const editorInit = () => {
  // 在下一个DOM更新周期执行
  nextTick(() => {
    // 配置JavaScript的语法检查选项
    monaco.languages.typescript.javascriptDefaults.setDiagnosticsOptions({
      noSemanticValidation: true,   // 禁用语义验证
      noSyntaxValidation: false     // 启用语法验证
    });
    
    // 配置JavaScript的编译器选项
    monaco.languages.typescript.javascriptDefaults.setCompilerOptions({
      target: monaco.languages.typescript.ScriptTarget.ES2016, // 编译目标ES2016
      allowNonTsExtensions: true    // 允许非TypeScript扩展
    });
    
    // 创建或重置编辑器
    if (!editor) {
      // 创建新的编辑器实例
      editor = monaco.editor.create(document.getElementById('codeEditBox') as HTMLElement, {
        value: text.value,          // 初始文本
        language: 'go',             // 默认语言
        automaticLayout: true,      // 自适应布局
        theme: 'vs-dark',           // 深色主题
        foldingStrategy: 'indentation', // 缩进折叠策略
        renderLineHighlight: 'all', // 高亮整行
        selectOnLineNumbers: true,  // 点击行号选择整行
        minimap: {                  // 缩略图设置
          enabled: false,           // 禁用缩略图
        },
        readOnly: false,            // 非只读模式
        fontSize: 16,               // 字体大小
        scrollBeyondLastLine: false, // 不在最后一行后滚动
        overviewRulerBorder: false, // 禁用概览标尺边框
      })
    } else {
      // 重置编辑器内容
      editor.setValue("");
    }
    
    // 监听编辑器内容变化
    editor.onDidChangeModelContent((val: any) => {
      text.value = editor.getValue(); // 更新代码文本
    })
  })
}

// 初始化编辑器
editorInit()

// 切换编程语言
// @ts-ignore
const changeLanguage = () => {
  // 设置编辑器的语言模型
  monaco.editor.setModelLanguage(editor.getModel(), language.value)
}

// 提交代码函数
const submitCode = () => {
  loading.value = true  // 显示加载状态
  // 调用API提交代码
  api.submitCode(text.value, route.query.identity).then(res => {
    loading.value = false  // 隐藏加载状态
    if (res.data.code == 200) {
      msg.value = res.data.data.msg  // 更新编译结果消息
      
      // 根据状态显示不同消息
      if (res.data.data.status == 1) {
        ElMessage.success(res.data.data.msg)  // 成功消息
      } else {
        ElMessage.warning(res.data.data.msg)  // 警告消息
      }
    } else {
      ElMessage.error(res.data.msg)  // 错误消息
    }
  })
}

/***
 * 编辑器常用方法示例（注释）：
 * 
 * editor.setValue(newValue)  // 设置编辑器内容
 * 
 * editor.getValue()          // 获取编辑器内容
 * 
 * editor.onDidChangeModelContent((val)=>{ //监听值的变化  })
 * 
 * editor.getAction('editor.action.formatDocument').run()    //格式化代码
 * 
 * editor.dispose()    //销毁实例
 * 
 * editor.onDidChangeOptions　　//配置改变事件
 * 
 * editor.onDidChangeLanguage　　//语言改变事件
 */
</script>

<style scoped lang="scss">
/* 编辑器容器样式 */
#codeEditBox {
  height: 50%;  /* 占据50%高度 */
}

/* 选择器样式 */
.select {
  padding: 10px 0;  /* 上下内边距 */
}

/* 提交按钮容器样式 */
.submit {
  text-align: center;  /* 居中 */
  padding: 10px 0;     /* 上下内边距 */
}

/* 编辑器整体容器 */
.e-box {
  width: 100%;         /* 100%宽度 */
  padding: 10px;       /* 内边距 */
  box-sizing: border-box; /* 包含padding在尺寸内 */
  height: 100%;        /* 100%高度 */
  display: flex;       /* 弹性布局 */
  justify-content: space-between; /* 两端对齐 */
  box-sizing: border-box;
  flex-direction: column; /* 垂直排列 */
}

/* 编译结果容器样式 */
.sub-box {
  border-radius: 10px;   /* 圆角 */
  background-color: #999; /* 灰色背景 */
  padding: 10px;         /* 内边距 */
  box-sizing: border-box; 
  flex: 1;               /* 占据剩余空间 */
  color: white;          /* 白色文字 */
}
</style>