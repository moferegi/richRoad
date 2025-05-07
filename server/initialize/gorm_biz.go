package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
)

func bizModel() error {
	db := global.GVA_DB
	err := db.AutoMigrate(client.ClientUser{}, shop.Banner{}, shop.Category{}, shop.Good{}, shop.Sku{}, shop.Cart{}, shop.Order{}, shop.OrderDetail{}, client.Address{}, client.Collect{}, shop.Comment{}, shop.Tag{})
	if err != nil {
		return err
	}
	return nil
}
