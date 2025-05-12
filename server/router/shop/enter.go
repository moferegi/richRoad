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
	TagRouter
	CouponRouter
	CouponOrderUserRouter
}

var (
	commentApi = api.ApiGroupApp.ShopApiGroup.CommentApi
	tagApi     = api.ApiGroupApp.ShopApiGroup.TagApi
	CouApi     = api.ApiGroupApp.ShopApiGroup.CouponApi
	couApi     = api.ApiGroupApp.ShopApiGroup.CouponOrderUserApi
)
