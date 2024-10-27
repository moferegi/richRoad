package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/router/client"
	"github.com/flipped-aurora/gin-vue-admin/server/router/example"
	"github.com/flipped-aurora/gin-vue-admin/server/router/shop"
	"github.com/flipped-aurora/gin-vue-admin/server/router/system"
)

var RouterGroupApp = new(RouterGroup)

type RouterGroup struct {
	System  system.RouterGroup
	Example example.RouterGroup
	Client  client.RouterGroup
	Shop    shop.RouterGroup
}
