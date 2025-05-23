package shop

import (
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
	"gorm.io/gorm"
	"time"
)

type OrderService struct {
}

// CreateOrder 创建订单记录
// Author [piexlmax](https://github.com/piexlmax)
func (orderService *OrderService) CreateOrder(order *shop.Order) (err error) {
	err = global.GVA_DB.Create(order).Error
	return err
}

func (orderService *OrderService) PlaceOrder(order *shop.Order) (OrderID uint, err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		order.TotalPrice = 0
		for i := range order.Detail {
			// 判断库存
			var sku shop.Sku
			err = tx.Where("id = ?", order.Detail[i].SKUID).First(&sku).Error
			if err != nil {
				return err
			}
			if sku.Inventory < order.Detail[i].Quantity {
				return errors.New("库存不足")
			}
			err = tx.Model(&shop.Sku{}).Where("id = ?", order.Detail[i].SKUID).Update("inventory", gorm.Expr("inventory - ?", order.Detail[i].Quantity)).Error
			if err != nil {
				return err
			}
			order.Detail[i].Price = sku.Price
			order.TotalPrice += order.Detail[i].Quantity * order.Detail[i].Price
			// 减扣库存
		}
		if order.CouponNum != "" {
			var couponOrderUser shop.CouponOrderUser
			err = tx.Where("coupon_num = ? and order_id IS NULL", order.CouponNum).First(&couponOrderUser).Error
			if err != nil {
				return err
			}
			var coupon shop.Coupon
			err = tx.Where("id = ?", couponOrderUser.CouponID).First(&coupon).Error
			if err != nil {
				return err
			}
			if coupon.ProductID != nil {
				hasGoods := false
				for _, detail := range order.Detail {
					if int(detail.ID) == *coupon.ProductID {
						hasGoods = true
						break
					}
				}
				if !hasGoods {
					return errors.New("当前订单不可使用此券")
				}
			}
			order.TotalPrice -= coupon.Discount
		}
		order.Status = "0"

		if addr, err := orderService.GetDefaultAddress(order.UserID); err == nil {
			order.Name = addr.Name
			order.Phone = addr.Phone
			order.Province = addr.ProvinceStr
			order.City = addr.CityStr
			order.Area = addr.AreaStr
			order.Street = addr.Street
		}
		order.CloseTime = time.Now().Add(time.Minute * 15)

		err = tx.Create(order).Error
		if err != nil {
			return err
		}
		OrderID = order.ID
		err = tx.Model(&shop.CouponOrderUser{}).Where("coupon_num = ?", order.CouponNum).Update("order_id", order.ID).Error
		if err != nil {
			return err
		}
		return nil
	})
	return
}

// 购物车下单
func (orderService *OrderService) PlaceOrderByCart(userID uint) (OrderID uint, err error) {
	// 1. 获取购物车中的商品
	var carts []shop.Cart
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		tx.Where("user_id = ?", userID).Preload("Good").Preload("SKU").Find(&carts)
		// 2. 创建订单
		order := shop.Order{
			UserID: userID,
			Status: "0",
		}
		order.TotalPrice = 0
		// 从购物车设置订单详情
		order.Detail = make([]shop.OrderDetail, 0)
		for i := range carts {
			// 判断当前库存是否充足
			if carts[i].SKU.Inventory < carts[i].Quantity {
				return errors.New("库存不足")
			}

			order.Detail = append(order.Detail, shop.OrderDetail{
				GoodID:   carts[i].GoodID,
				SKUID:    carts[i].SKUID,
				Quantity: carts[i].Quantity,
				Price:    carts[i].SKU.Price,
			})
			order.TotalPrice += carts[i].Quantity * carts[i].SKU.Price

			// 扣减库存 增加销量
			err = tx.Model(&shop.Sku{}).Where("id = ?", carts[i].SKUID).
				Update("inventory", gorm.Expr("inventory - ?", carts[i].Quantity)).
				Update("sale_num", gorm.Expr("sale_num + ?", carts[i].Quantity)).
				Error
			if err != nil {
				return err
			}
			err = tx.Model(&shop.Good{}).Where("id = ?", carts[i].GoodID).Update("sale_num", gorm.Expr("sale_num + ?", carts[i].Quantity)).Error
			if err != nil {
				return err
			}
		}

		if addr, err := orderService.GetDefaultAddress(order.UserID); err == nil {
			order.Name = addr.Name
			order.Phone = addr.Phone
			order.Province = addr.ProvinceStr
			order.City = addr.CityStr
			order.Area = addr.AreaStr
			order.Street = addr.Street
		}
		order.CloseTime = time.Now().Add(time.Minute * 15)

		// 3. 创建订单详情

		err = tx.Create(&order).Error
		if err != nil {
			return err
		}

		// 4. 删除购物车中的商品
		err = tx.Where("user_id = ?", userID).Delete(&shop.Cart{}).Error
		if err != nil {
			return err
		}
		OrderID = order.ID
		return nil
	})
	return
}

func (orderService *OrderService) UpdateOrderStatus(db *gorm.DB, orderID string, status string) (err error) {
	if db == nil {
		db = global.GVA_DB
	}
	err = db.Transaction(func(tx *gorm.DB) error {
		err = tx.Model(&shop.Order{}).Where("id = ?", orderID).Update("status", status).Error
		if err != nil {
			return err
		}
		if status == "5" {
			var order shop.Order
			err = tx.Preload("Detail").Where("id = ?", orderID).First(&order).Error
			if err != nil {
				return err
			}
			for _, detail := range order.Detail {
				err = tx.Model(&shop.Sku{}).Where("id = ?", detail.SKUID).Update("inventory", gorm.Expr("inventory + ?", detail.Quantity)).Error
				if err != nil {
					return err
				}
			}
			return nil
		}
		return err
	})
	return err
}

// DeleteOrder 删除订单记录
// Author [piexlmax](https://github.com/piexlmax)
func (orderService *OrderService) DeleteOrder(ID string) (err error) {
	err = global.GVA_DB.Delete(&shop.Order{}, "id = ?", ID).Error
	return err
}

// DeleteOrderByIds 批量删除订单记录
// Author [piexlmax](https://github.com/piexlmax)
func (orderService *OrderService) DeleteOrderByIds(IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]shop.Order{}, "id in ?", IDs).Error
	return err
}

// UpdateOrder 更新订单记录
// Author [piexlmax](https://github.com/piexlmax)
func (orderService *OrderService) UpdateOrder(order shop.Order) (err error) {
	err = global.GVA_DB.Model(&shop.Order{}).Where("id = ?", order.ID).Updates(&order).Error
	return err
}

// GetOrder 根据ID获取订单记录
// Author [piexlmax](https://github.com/piexlmax)
func (orderService *OrderService) GetOrder(ID string, userID uint) (order shop.OrderRes, err error) {
	db := global.GVA_DB.Where("id = ?", ID).Preload("Detail").Preload("Detail.Good").Preload("Detail.SKU").Preload("Comment")
	if userID != 0 {
		db = db.Where("user_id = ?", userID)
	}
	err = db.First(&order).Error
	return
}

func (orderService *OrderService) SelfOrderComment(ID, goodID, SKUID string, userID uint) (order shop.OrderCommentRes, err error) {
	var comment shop.Comment
	var detail shop.OrderDetailRes
	_ = global.GVA_DB.First(&comment, "order_id = ? && good_id = ? && sku_id = ?", ID, goodID, SKUID).Error

	err = global.GVA_DB.Preload("Good").Preload("SKU").First(&detail, "order_id = ? && good_id = ? && sku_id = ?", ID, goodID, SKUID).Error
	if err != nil {
		return order, err
	}
	db := global.GVA_DB.Where("id = ?", ID)
	if userID != 0 {
		db = db.Where("user_id = ?", userID)
	}
	err = db.First(&order).Error
	order.Comment = comment
	order.Detail = detail
	return
}

// GetOrderInfoList 分页获取订单记录
// Author [piexlmax](https://github.com/piexlmax)
func (orderService *OrderService) GetOrderInfoList(info shopReq.OrderSearch) (list []shop.OrderRes, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	// 创建db
	db := global.GVA_DB.Model(&shop.OrderRes{}).Preload("Detail").Preload("Detail.Good").Preload("Detail.SKU")
	var orders []shop.OrderRes

	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.UserID != nil {
		db = db.Where("user_id = ?", info.UserID)
	}

	// 根据info.status查询
	if info.Status != "" {
		db = db.Where("status = ?", info.Status)
	}

	cErr := db.Count(&total).Error
	if cErr != nil {
		return
	}
	// 应用分页
	if limit > 0 {
		db = db.Limit(limit).Offset(offset)
	}

	// 执行查询
	err = db.Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// GetDefaultAddress 获取用户默认地址 如果用户没有设置则拉取表内最后一条
func (orderService *OrderService) GetDefaultAddress(UserID uint) (address client.Address, err error) {
	// 尝试获取默认地址
	err = global.GVA_DB.Where("user_id = ? and active = ?", UserID, true).First(&address).Error
	if err != nil {
		// 如果所有active都为0，则使用Order()按表内所有时间排序后返回最后一条
		err = global.GVA_DB.Where("user_id = ?", UserID).Order("created_at desc").First(&address).Error
	}
	return address, nil
}
