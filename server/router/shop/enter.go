package shop

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct {
	BannerRouter
	CategoryRouter
	GoodRouter
	SkuRouter
	CartRouter
	OrderRouter
	CommentRouter
}

var commentApi = api.ApiGroupApp.ShopApiGroup.CommentApi
