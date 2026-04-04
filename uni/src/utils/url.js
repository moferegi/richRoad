import { baseUrl } from "@/utils/request.js"
import { getDefaultDomain } from "@/api/sysConfig.js"

// 缓存外部域名
let cachedDomain = ""
let domainLoading = false

export const initExternalDomain = async () => {
	if (cachedDomain || domainLoading) return
	domainLoading = true
	try {
		const res = await getDefaultDomain()
		if (res.code === 0 && res.data) {
			cachedDomain = res.data.replace(/\/+$/, "")
		}
	} catch (e) {
		console.warn("获取外部域名失败", e)
	} finally {
		domainLoading = false
	}
}

export const getUrl = (url) => {
	if (!url) return "http://localhost:8888"
	const path = baseUrl
	if (url.slice(0, 4) === 'data') {
		return url
	}
	if (url && url.slice(0, 4) !== 'http') {
		if (path === "/") {
			return url
		}
		if (url.slice(0, 1) === "/") {
			return path + url
		}
		return path + "/" + url
	} else {
		return url
	}
}

// 获取外部链接完整URL：如果是相对路径，自动拼接外部默认域名
export const getExternalUrl = (url) => {
	if (!url) return ""
	if (url.slice(0, 4) === 'http') return url
	if (url.slice(0, 4) === 'data') return url
	if (cachedDomain) {
		const sep = url.slice(0, 1) === "/" ? "" : "/"
		return cachedDomain + sep + url
	}
	// 兜底使用 getUrl
	return getUrl(url)
}