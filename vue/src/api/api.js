import http from './http.js' // 导入自定义HTTP模块

// API方法集合
export default {
	// 获取问题列表
	getProblemList(param) {
		return http.get(`/problem-list`, param)
	},
	
	// 获取问题详情
	getProblemDetail(param) {
		return http.get(`/problem-detail`, param)
	},
	
	// 获取分类列表
	getSortList(param) {
		return http.get(`/category-list`, param)
	},
	
	// 获取排行榜数据
	getRankList(param) {
		return http.get(`/rank-list`, param)
	},
	
	// 获取提交列表
	getSubmitList(param) {
		return http.get(`/submit-list`, param)
	},
	
	// 发送验证码
	sendCode(param) {
		return http.postUncode(`/send-code`, param)
	},
	
	// 用户登录
	login(param) {
		return http.postUncode(`/login`, param)
	},
	
	// 用户注册
	register(param) {
		return http.postUncode(`/register`, param)
	},
	
	// 删除分类（管理员）
	delSort(param) {
		return http.delete(`/admin/category-delete`, param)
	},
	
	// 创建分类（管理员）
	addSort(param) {
		return http.postUncode(`/admin/category-create`, param)
	},
	
	// 创建问题（管理员）
	addProblem(param) {
		return http.post(`/admin/problem-create`, param)
	},
	
	// 编辑问题（管理员）
	editProblem(param) {
		return http.putJson(`/admin/problem-modify`, param)
	},
	
	// 编辑分类（管理员）
	editSort(param) {
		return http.put(`/admin/category-modify`, param)
	},
	
	// 提交代码
	submitCode(param, id) {
		return http.postJson(`/user/submit?problem_identity=${id}`, param)
	},
	
	// 文件上传
	uploadFile(param) {
		return http.upFile('', param)
	}
}