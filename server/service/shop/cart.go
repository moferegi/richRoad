package shop

import (
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
)

type CartService struct {
}

// CreateCart 创建购物车记录
// Author [piexlmax](https://github.com/piexlmax)
func (cartService *CartService) CreateCart(cart *shop.Cart) (err error) {
	err = global.GVA_DB.Create(cart).Error
	return err
}

// DeleteCart 删除购物车记录
// Author [piexlmax](https://github.com/piexlmax)
func (cartService *CartService) DeleteCart(ID string) (err error) {
	err = global.GVA_DB.Delete(&shop.Cart{}, "id = ?", ID).Error
	return err
}

// DeleteCartByIds 批量删除购物车记录
// Author [piexlmax](https://github.com/piexlmax)
func (cartService *CartService) DeleteCartByIds(IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]shop.Cart{}, "id in ? ", IDs).Error
	return err
}

// UpdateCart 更新购物车记录
// Author [piexlmax](https://github.com/piexlmax)
func (cartService *CartService) UpdateCart(cart shop.Cart) (err error) {
	err = global.GVA_DB.Model(&shop.Cart{}).Where("id = ?", cart.ID).Updates(&cart).Error
	return err
}

// GetCart 根据ID获取购物车记录
// Author [piexlmax](https://github.com/piexlmax)
func (cartService *CartService) GetCart(ID string) (cart shop.Cart, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&cart).Error
	return
}

// GetCartInfoList 分页获取购物车记录
// Author [piexlmax](https://github.com/piexlmax)
func (cartService *CartService) GetCartInfoList(info shopReq.CartSearch) (list []shop.Cart, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&shop.Cart{})
	var carts []shop.Cart
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}

	if info.UserID != 0 {
		db = db.Where("user_id = ?", info.UserID)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&carts).Error
	return carts, total, err
}

func (cartService *CartService) GetSelfCart(userID uint) (list []shop.Cart, err error) {
	err = global.GVA_DB.Preload("Good").Preload("SKU").Find(&list, "user_id = ?", userID).Error
	return
}

func (cartService *CartService) AddCart(cartReq *shopReq.CartCreate) (err error) {
	var cart shop.Cart
	if cartReq.Quantity == 0 {
		cartReq.Quantity = 1
	}
	var sku shop.Sku
	if err = global.GVA_DB.First(&sku, "id = ?", cartReq.SKUID).Error; err != nil {
		return errors.New("orderGoodUnavailable")
	}
	if sku.GoodID != cartReq.GoodID {
		return errors.New("orderGoodUnavailable")
	}
	if _, err = ensurePublicGoodPurchasable(global.GVA_DB, sku.GoodID); err != nil {
		return err
	}
	if sku.Inventory < cartReq.Quantity {
		return errors.New("orderInventoryInsufficient")
	}

	ferr := global.GVA_DB.First(&cart, "user_id = ? AND good_id = ? AND sku_id = ?", cartReq.UserID, cartReq.GoodID, cartReq.SKUID).Error
	if ferr != nil {
		cart.Quantity = cartReq.Quantity
		cart.UserID = cartReq.UserID
		cart.GoodID = cartReq.GoodID
		cart.SKUID = cartReq.SKUID
		err = global.GVA_DB.Create(&cart).Error
		return
	}

	quantity := cart.Quantity + cartReq.Quantity
	err = global.GVA_DB.Model(&cart).Update("quantity", quantity).Error
	return
}

func (cartService *CartService) CutCart(cartReq *shopReq.CartCreate) (err error) {
	var cart shop.Cart
	ferr := global.GVA_DB.First(&cart, "user_id = ? AND good_id = ? AND sku_id = ?", cartReq.UserID, cartReq.GoodID, cartReq.SKUID).Error
	if ferr != nil {
		return errors.New("购物车中没有该商品")
	}
	quantity := cart.Quantity - 1
	if quantity == 0 {
		err = global.GVA_DB.Delete(&cart).Error
		return
	}
	err = global.GVA_DB.Model(&cart).Update("quantity", quantity).Error
	return

}

func (cartService *CartService) ClearCart(userID uint) (err error) {
	err = global.GVA_DB.Delete(&shop.Cart{}, "user_id = ?", userID).Error
	return
}
