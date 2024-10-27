package task

import (
	"errors"
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
