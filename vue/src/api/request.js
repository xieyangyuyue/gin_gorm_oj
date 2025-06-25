/****   request.js   ****/
// 导入axios库和qs序列化库
import axios from 'axios'
import qs from 'qs'
// 导入Element Plus的消息提示组件
import { ElMessage } from 'element-plus'

// 创建axios实例
const service = axios.create({
    // 基础URL，开发环境使用本地服务
		// baseURL: 'http://111.170.11.164:10700/',
	baseURL: 'http://localhost:8080/',
	// 超时时间设置为300秒（5分钟）
	timeout: 300 * 1000
})

// 请求拦截器
service.interceptors.request.use(config => {
    // 处理GET请求的参数序列化
		if (config.method === 'get') {
			config.paramsSerializer = function(params) {
                // 使用qs序列化参数，数组格式为逗号分隔
				return qs.stringify(params, {
					arrayFormat: 'comma'
				})
			}
		}
		
        // 从localStorage获取token
		const token=localStorage.token
        // 如果token存在，添加到请求头
		if(token){
			config.headers['Authorization']=token
		}
	return config
}, error => {
    // 请求错误处理
	Promise.reject(error)
})

// 响应拦截器
service.interceptors.response.use((config) => {
    // 正常响应直接返回
	return config
}, (error) => {
    // 错误响应处理
	if (error.response) {
        // 获取错误信息
		const errorMessage = error.response.data === null 
            ? '系统内部异常，请联系网站管理员' 
            : error.response.data.message
            
        // 根据HTTP状态码处理不同错误
		switch (error.response.status) {
			case 404: // 资源未找到
				ElMessage('很抱歉，资源未找到')
				break
			case 403: // 无操作权限
				ElMessage('很抱歉，您暂无该操作权限')
				break
			case 401: // 认证失效
				ElMessage('很抱歉，认证已失效，请重新登录')
				// let aid=getQueryVariable('corpId')
				// localStorage.removeItem(aid)
				// alert(aid)
				break
			default: // 其他错误
				if (errorMessage === 'refresh token无效') {
					ElMessage('登录已过期，请重新登录')
				} else {
					ElMessage(errorMessage)
				}
				break
		}
	}
	return Promise.reject(error)
})

// 导出配置好的axios实例
export default service