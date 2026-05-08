package initialize

import (
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	wxpayModel "github.com/flipped-aurora/gin-vue-admin/server/plugin/wxpay/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const migrateClientTryonPointKey = "migration_client_user_tryon_point_v1"

func bizModel() error {
	db := global.GVA_DB
	err := db.AutoMigrate(client.ClientUser{}, shop.Banner{}, shop.Category{}, shop.Good{}, shop.Sku{}, shop.Cart{}, shop.CreateOrder{}, shop.OrderDetail{}, client.Address{}, client.Collect{}, shop.Comment{}, shop.Tag{}, shop.Coupon{}, shop.CouponOrderUser{}, shop.Promotion{}, shop.History{}, client.PointRecord{}, client.TryonRechargeOrder{}, client.TryonTask{}, client.TryonModel{}, wxpayModel.Order{}, shop.Kefu{}, client.VisitorLog{}, client.VisitorSummary{}, client.SysConfig{}, shop.GoodPurchase{}, shop.QrcodePayment{}, shop.Popup{}, shop.MarketingReward{}, client.SysLanguage{}, client.PhoneAreaCode{}, client.SignIn{}, client.ExternalLinkDomain{}, client.LoginFailRecord{}, client.RegisterIPRecord{}, shop.SkuSpec{})
	if err != nil {
		return err
	}
	if err = migrateClientTryonPoint(db); err != nil {
		return err
	}
	return nil
}
func migrateClientTryonPoint(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var marker client.SysConfig
		err := tx.Where("config_key = ?", migrateClientTryonPointKey).First(&marker).Error
		if err == nil {
			return nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if err = tx.Model(&client.ClientUser{}).Where("point > 0 AND (tryon_point = 0 OR tryon_point IS NULL)").Update("tryon_point", gorm.Expr("point")).Error; err != nil {
			return err
		}
		flag := client.SysConfig{ConfigKey: migrateClientTryonPointKey, ConfigValue: "1", ConfigName: "用户资产拆分迁移标记", ConfigGroup: "migration", Remark: "v1: 初始化 tryon_point=point，避免历史资产丢失"}
		return tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&flag).Error
	})
}
