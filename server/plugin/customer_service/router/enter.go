package router

// RouterGroup 路由组
type RouterGroup struct {
	CustomerServiceRouter
}

// Router 全局路由实例
var Router = new(RouterGroup)
