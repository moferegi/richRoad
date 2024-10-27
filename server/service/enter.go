package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/service/client"
	"github.com/flipped-aurora/gin-vue-admin/server/service/example"
	"github.com/flipped-aurora/gin-vue-admin/server/service/shop"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
)

var ServiceGroupApp = new(ServiceGroup)

type ServiceGroup struct {
	SystemServiceGroup  system.ServiceGroup
	ExampleServiceGroup example.ServiceGroup
	ClientServiceGroup  client.ServiceGroup
	ShopServiceGroup    shop.ServiceGroup
}
