package task

import (
	"errors"
	"fmt"
	"time"

	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"gorm.io/gorm"
)

func ClearOrder(db *gorm.DB) error {
	if db == nil {
		return errors.New("db Cannot be empty")
	}

	var orders []shop.Order
	now := time.Now()
	err := db.Where("status IN ? AND close_time < ?", []string{"0", "8"}, now).Limit(500).Find(&orders).Error
	if err != nil {
		return err
	}
	var lastErr error
	for _, order := range orders {
		e := service.ServiceGroupApp.ShopServiceGroup.UpdateOrderStatus(db, strconv.Itoa(int(order.ID)), "4")
		if e != nil {
			fmt.Printf("[ClearOrder] 自动取消订单 %d 失败: %v\n", order.ID, e)
			lastErr = e
			continue
		}
		fmt.Printf("[ClearOrder] 自动取消订单 %d 成功\n", order.ID)
	}
	return lastErr
}

// ClearExpiredOrders 清理超时未支付的订单（基于expire_at字段）
func ClearExpiredOrders(db *gorm.DB) error {
	if db == nil {
		return errors.New("db Cannot be empty")
	}

	now := time.Now()
	var orders []shop.Order
	// 查找待支付且已过期的订单（每批最多500条）
	err := db.Where("status = ? AND expire_at IS NOT NULL AND expire_at < ?", "0", now).Limit(500).Find(&orders).Error
	if err != nil {
		return err
	}
	var lastErr error
	for _, order := range orders {
		e := service.ServiceGroupApp.ShopServiceGroup.UpdateOrderStatus(db, strconv.Itoa(int(order.ID)), "4")
		if e != nil {
			fmt.Printf("[ClearExpiredOrders] 自动取消订单 %d 失败: %v\n", order.ID, e)
			lastErr = e
			continue
		}
		fmt.Printf("[ClearExpiredOrders] 自动取消订单 %d 成功\n", order.ID)
	}
	return lastErr
}

// ClearExpiredTryonRechargeOrders 清理超时未支付的试衣币充值订单
func ClearExpiredTryonRechargeOrders(db *gorm.DB) error {
	if db == nil {
		return errors.New("db Cannot be empty")
	}

	now := time.Now()
	var orders []client.TryonRechargeOrder
	err := db.Where("status IN ? AND close_time < ?", []string{"0", "8"}, now).Limit(500).Find(&orders).Error
	if err != nil {
		return err
	}

	var lastErr error
	for _, order := range orders {
		e := db.Model(&client.TryonRechargeOrder{}).Where("id = ? AND status IN ?", order.ID, []string{"0", "8"}).Updates(map[string]interface{}{
			"status":       "4",
			"cancelled_at": now,
		}).Error
		if e != nil {
			fmt.Printf("[ClearExpiredTryonRechargeOrders] 自动取消充值订单 %d 失败: %v\n", order.ID, e)
			lastErr = e
			continue
		}
		fmt.Printf("[ClearExpiredTryonRechargeOrders] 自动取消充值订单 %d 成功\n", order.ID)
	}
	return lastErr
}
