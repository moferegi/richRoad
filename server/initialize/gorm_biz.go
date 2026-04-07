package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	wxpayModel "github.com/flipped-aurora/gin-vue-admin/server/plugin/wxpay/model"
)

func bizModel() error {
	db := global.GVA_DB
	err := db.AutoMigrate(client.ClientUser{}, shop.Banner{}, shop.Category{}, shop.Good{}, shop.Sku{}, shop.Cart{}, shop.CreateOrder{}, shop.OrderDetail{}, client.Address{}, client.Collect{}, shop.Comment{}, shop.Tag{}, shop.Coupon{}, shop.CouponOrderUser{}, shop.Promotion{}, shop.History{}, client.PointRecord{}, wxpayModel.Order{}, shop.Kefu{}, client.VisitorLog{}, client.VisitorSummary{}, client.SysConfig{}, shop.GoodPurchase{}, shop.QrcodePayment{}, shop.Popup{}, shop.MarketingReward{}, client.SysLanguage{}, client.PhoneAreaCode{}, client.SignIn{}, client.ExternalLinkDomain{}, client.LoginFailRecord{}, client.RegisterIPRecord{}, shop.SkuSpec{})
	if err != nil {
		return err
	}
	return nil
}
