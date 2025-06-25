<template>
  <!-- Monaco 编辑器容器 -->
  <div>
    <!-- 编辑器挂载点 -->
    <div id="monaco" ref="mymonaco"></div>
  </div>
</template>

<script lang="ts" setup>
// 导入Vue生命周期钩子
import { onMounted } from '@vue/runtime-core'

// 导入Monaco Editor API（核心功能）
import * as monaco from 'monaco-editor/esm/vs/editor/editor.api.js'

// 导入JavaScript语言支持（增强JavaScript语法功能）
import 'monaco-editor/esm/vs/basic-languages/javascript/javascript.contribution'

// 导入Vue响应式API
import { ref } from 'vue'

// 创建编辑器容器的引用
const mymonaco = ref<HTMLElement | null>()

// 组件挂载后的生命周期钩子
onMounted(
  () => {
    // 以下是被注释掉的worker配置（用于处理语言服务）
    // @ts-ignore
    // self.MonacoEnvironment = {
    //   getWorker: function (_: any, label: string) {
    //     console.log(label)
    //     if (label === 'json') {
    //       return new jsonWorker()
    //     }
    //     new editorWorker()
    //   },
    // }

    // 检查编辑器容器元素是否存在
    if (mymonaco.value) {
      // 创建Monaco编辑器实例
      const monacoInstance = monaco.editor.create(mymonaco.value, {
        value: 'console.log("aaa")',  // 初始代码内容
        language: 'javascript'       // 默认语言为JavaScript
      })
      // 销毁编辑器实例（当前注释状态）
      //  monacoInstance.dispose()
    }
  }
)
</script>

<style scoped lang="scss">
/* Monaco编辑器容器样式 */
#monaco {
  width: 500px;   /* 固定宽度 */
  height: 500px;  /* 固定高度 */
}
</style>