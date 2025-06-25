/****   http.js   ****/
// 导入封装好的axios实例和qs序列化库
import request from './request'
import qs from 'qs'

// HTTP方法封装对象
const http = {
    /**
     * GET请求方法
     * @param {string} url 请求地址 
     * @param {Object} params 请求参数
     */
    get(url, params) {
        // 配置请求对象
        const config = {
            method: 'get',
            url: url
        }
        // 如果有参数则添加
        if (params) config.params = params
        return request(config)
    },

    /**
     * POST请求方法（默认JSON格式）
     * @param {string} url 请求地址
     * @param {Object} params 请求参数
     */
    post(url, params) {
        const config = {
            method: 'post',
            url: url
        }
        if (params) config.data = params
        return request(config)
    },

    /**
     * POST请求方法（multipart/form-data格式）
     * @param {string} url 请求地址
     * @param {Object} params 请求参数
     */
    postJson(url, params) {
        const config = {
            method: 'post',
            url: url,
            headers: {
                'Content-Type': 'multipart/form-data'
            }
        }
        console.log(params) // 打印参数（调试用）
        if (params) config.data = params
        return request(config)
    },

    /**
     * POST请求方法（x-www-form-urlencoded格式，支持数组参数）
     * @param {string} url 请求地址
     * @param {Object} params 请求参数
     */
    postSB(url, params) {
        const config = {
            method: 'post',
            url: url,
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded'
            }
        }
        // 使用qs序列化参数（数组格式为repeat）
        if (params) config.data = qs.stringify(params, { arrayFormat: 'repeat' })
        return request(config)
    },

    /**
     * POST请求方法（x-www-form-urlencoded格式）
     * @param {string} url 请求地址
     * @param {Object} params 请求参数
     */
    postUncode(url, params) {
        const config = {
            method: 'post',
            url: url,
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded'
            }
        }
        // 使用qs序列化参数
        if (params) config.data = qs.stringify(params)
        return request(config)
    },

    /**
     * 文件上传方法
     * @param {string} url 请求地址
     * @param {Object} data 文件数据
     */
    upFile(url, data) {
        return request({
            method: 'post',
            url: url,
            data: data,
            headers: {
                "Content-Type": "multipart/form-data"
            }
        })
    },

    /**
     * PUT请求方法（x-www-form-urlencoded格式）
     * @param {string} url 请求地址
     * @param {Object} params 请求参数
     */
    put(url, params) {
        const config = {
            method: 'put',
            url: url,
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded'
            }
        }
        // 使用qs序列化参数
        if (params) config.data = qs.stringify(params)
        return request(config)
    },

    /**
     * PUT请求方法（JSON格式）
     * @param {string} url 请求地址
     * @param {Object} params 请求参数
     */
    putJson(url, params) {
        const config = {
            method: 'put',
            url: url,
        }
        if (params) config.data = params
        return request(config)
    },

    /**
     * DELETE请求方法
     * @param {string} url 请求地址
     * @param {Object} params 请求参数
     */
    delete(url, params) {
        const config = {
            method: 'delete',
            url: url,
        }
        if (params) config.params = params
        return request(config)
    }
}

// 导出HTTP对象
export default http