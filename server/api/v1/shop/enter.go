package shop

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	BannerApi
	CategoryApi
	GoodApi
	SkuApi
	CartApi
	OrderApi
	CommentApi
	TagApi
	CouponApi
	CouponOrderUserApi
	PromotionApi
}

var (
	commentService = service.ServiceGroupApp.ShopServiceGroup.CommentService
	tagService     = service.ServiceGroupApp.ShopServiceGroup.TagService
	couponService  = service.ServiceGroupApp.ShopServiceGroup.CouponService
	couService     = service.ServiceGroupApp.ShopServiceGroup.CouponOrderUserService
	PromoService   = service.ServiceGroupApp.ShopServiceGroup.PromotionService
)
