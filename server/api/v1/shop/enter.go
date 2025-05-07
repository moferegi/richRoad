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
}

var (
	commentService = service.ServiceGroupApp.ShopServiceGroup.CommentService
	tagService     = service.ServiceGroupApp.ShopServiceGroup.TagService
)
