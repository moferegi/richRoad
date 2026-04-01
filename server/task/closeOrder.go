package task

import (
	"errors"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"gorm.io/gorm"
	"strconv"
)

func ClearOrder(db *gorm.DB) error {
	if db == nil {
		return errors.New("db Cannot be empty")
	}

	var orders []shop.Order
	err := db.Where("status = ? and close_time < ?", "0", "now()").Find(&orders).Error
	if err != nil {
		return err
	}
	for _, order := range orders {
		e := service.ServiceGroupApp.ShopServiceGroup.UpdateOrderStatus(db, strconv.Itoa(int(order.ID)), "5")
		if e != nil {
			return e
		}
	}
	return nil
}

// ClearExpiredOrders 清理超时未支付的订单（基于expire_at字段）
func ClearExpiredOrders(db *gorm.DB) error {
	if db == nil {
		return errors.New("db Cannot be empty")
	}

	now := time.Now()
	var orders []shop.Order
	// 查找待支付且已过期的订单
	err := db.Where("status = ? AND expire_at IS NOT NULL AND expire_at < ?", "0", now).Find(&orders).Error
	if err != nil {
		return err
	}
	for _, order := range orders {
		e := service.ServiceGroupApp.ShopServiceGroup.UpdateOrderStatus(db, strconv.Itoa(int(order.ID)), "5")
		if e != nil {
			return e
		}
	}
	return nil
}
