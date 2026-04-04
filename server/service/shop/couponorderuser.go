package shop

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system" // 假设用户模型
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CouponOrderUserService struct{}

// 雪花算法节点，用于生成唯一ID
var (
	snowflakeNode *snowflake.Node
	nodeOnce      sync.Once
	nodeErr       error
)

// generateSnowflakeID 生成雪花算法ID
func generateSnowflakeID() (string, error) {
	nodeOnce.Do(func() {
		// 创建一个节点，节点ID为1
		snowflakeNode, nodeErr = snowflake.NewNode(1)
	})

	if nodeErr != nil {
		return "", nodeErr
	}

	// 生成ID
	id := snowflakeNode.Generate()
	return id.String(), nil
}

// CreateCouponOrderUser 创建优惠券记录
// Author [yourname](https://github.com/yourname)
func (couService *CouponOrderUserService) CreateCouponOrderUser(ctx context.Context, cou *shop.CouponOrderUser) (err error) {
	err = global.GVA_DB.Create(cou).Error
	return err
}

// DeleteCouponOrderUser 删除优惠券记录
// Author [yourname](https://github.com/yourname)
func (couService *CouponOrderUserService) DeleteCouponOrderUser(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&shop.CouponOrderUser{}, "id = ?", ID).Error
	return err
}

// DeleteCouponOrderUserByIds 批量删除优惠券记录
// Author [yourname](https://github.com/yourname)
func (couService *CouponOrderUserService) DeleteCouponOrderUserByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]shop.CouponOrderUser{}, "id in ?", IDs).Error
	return err
}

// UpdateCouponOrderUser 更新优惠券记录
// Author [yourname](https://github.com/yourname)
func (couService *CouponOrderUserService) UpdateCouponOrderUser(ctx context.Context, cou shop.CouponOrderUser) (err error) {
	err = global.GVA_DB.Model(&shop.CouponOrderUser{}).Where("id = ?", cou.ID).Updates(&cou).Error
	return err
}

// GetCouponOrderUser 根据ID获取优惠券记录
// Author [yourname](https://github.com/yourname)
func (couService *CouponOrderUserService) GetCouponOrderUser(ctx context.Context, ID string) (cou shop.CouponOrderUser, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&cou).Error
	return
}

// GetCouponOrderUserInfoList 分页获取优惠券记录
// Author [yourname](https://github.com/yourname)
func (couService *CouponOrderUserService) GetCouponOrderUserInfoList(ctx context.Context, info shopReq.CouponOrderUserSearch) (list []shop.CouponOrderUser, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&shop.CouponOrderUser{})
	var cous []shop.CouponOrderUser
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}

	if info.CouponID != nil {
		db = db.Where("coupon_id = ?", *info.CouponID)
	}
	if info.OrderID != nil {
		db = db.Where("order_id = ?", *info.OrderID)
	}
	if info.Status != nil {
		db = db.Where("status = ?", *info.Status)
	}
	if info.UserID != nil {
		db = db.Where("user_id = ?", *info.UserID)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&cous).Error
	return cous, total, err
}
func (couService *CouponOrderUserService) GetCouponOrderUserDataSource(ctx context.Context) (res map[string][]map[string]any, err error) {
	res = make(map[string][]map[string]any)

	couponID := make([]map[string]any, 0)

	global.GVA_DB.Table("shop_coupon").Where("deleted_at IS NULL").Select("name as label,id as value").Scan(&couponID)
	res["couponID"] = couponID
	orderID := make([]map[string]any, 0)

	global.GVA_DB.Table("shop_order").Where("deleted_at IS NULL").Select("out_trade_no as label,id as value").Scan(&orderID)
	res["orderID"] = orderID
	return
}
func (couService *CouponOrderUserService) GetCouponOrderUserPublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}

func (couService *CouponOrderUserService) GetAllClaimCoupon(ctx context.Context, userID uint, goodIds []int) (coupons []map[string]interface{}, err error) {
	// 初始化返回结果切片
	coupons = make([]map[string]interface{}, 0)

	// 1. 查询所有可用的优惠券
	var availableCoupons []shop.Coupon
	db := global.GVA_DB.Model(&shop.Coupon{})

	// 查询条件：优惠券有效且在有效期内
	db = db.Where("status = ? AND start_time <= NOW() AND end_time >= NOW()", true)

	// 执行查询
	if err = db.Find(&availableCoupons).Error; err != nil {
		return nil, errors.New("查询优惠券失败: " + err.Error())
	}

	// 2. 遍历优惠券，查询用户是否已领取
	for _, coupon := range availableCoupons {
		// 创建返回数据结构
		couponData := map[string]interface{}{
			"couponNum":       0, // 默认为0，表示未领取
			"couponID":        coupon.ID,
			"name":            *coupon.Name,
			"minSpend":        coupon.MinSpend,
			"discount":        coupon.Discount,
			"productID":       0,
			"startTime":       coupon.StartTime.Format("2006-01-02"),
			"endTime":         coupon.EndTime.Format("2006-01-02"),
			"status":          0, // 默认为0，表示未使用
			"canUse":          0,
			"nameI18n":        coupon.NameI18n,
			"descriptionI18n": coupon.DescriptionI18n,
			"backgroundImage": coupon.BackgroundImage,
			"externalBgPath":  coupon.ExternalBgPath,
		}

		// 如果有关联商品，设置productID
		if coupon.ProductID != nil {
			couponData["productID"] = *coupon.ProductID
		}
		if coupon.Description != nil {
			couponData["description"] = *coupon.Description
		}

		if len(goodIds) > 0 {
			// 检查商品ID是否在goodIds中 — 支持 ProductIDs (逗号分隔多商品)
			allowedIDs := couponAllowedProductIDs(coupon)
			if len(allowedIDs) == 0 {
				couponData["canUse"] = 1 // 没有商品限制，默认可用
			} else {
				for _, goodID := range goodIds {
					if allowedIDs[goodID] {
						couponData["canUse"] = 1
						break
					}
				}
			}
		} else {
			couponData["canUse"] = 1 // 没有商品限制，默认可用
		}

		// 查询用户是否已领取此优惠券
		var userCoupon shop.CouponOrderUser
		shopUserID := int(userID)
		result := global.GVA_DB.Where("coupon_id = ? AND shop_user_id = ?", coupon.ID, shopUserID).First(&userCoupon)

		// 如果用户已领取，设置券码和使用状态
		if result.Error == nil {
			couponData["couponNum"] = userCoupon.CouponNum
			if userCoupon.OrderID != nil {
				couponData["status"] = 1 // 占用状态 不能领取 不能使用
			}
		}

		// 添加到结果集
		coupons = append(coupons, couponData)
	}

	return coupons, nil
}

// ClaimCouponByUser 用户领取优惠券
func (couService *CouponOrderUserService) ClaimCouponByUser(ctx context.Context, userID uint, couponID int) (couponNum string, err error) {
	// 1. 检查优惠券是否存在、有效、库存充足
	var coupon shop.Coupon
	err = global.GVA_DB.Where("id = ? AND status = ? AND start_time <= NOW() AND end_time >= NOW() AND quantity > claimed", couponID, true).First(&coupon).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", errors.New("优惠券无效或已领完")
		}
		return "", errors.New("查询优惠券失败: " + err.Error())
	}

	// 2. 检查用户是否已领取该优惠券
	var existingClaim shop.CouponOrderUser
	shopUserID := int(userID) // Convert uint to int for ShopUserID
	err = global.GVA_DB.Where("coupon_id = ? AND shop_user_id = ?", couponID, shopUserID).First(&existingClaim).Error
	if err == nil {
		return "", errors.New("您已领取过该优惠券")
	}
	if err != gorm.ErrRecordNotFound {
		return "", errors.New("检查领取状态失败: " + err.Error())
	}

	// 3. 创建优惠券领取记录并更新优惠券已领取数量 (使用事务保证原子性)
	tx := global.GVA_DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err = tx.Error; err != nil {
		return "", err
	}

	// 生成雪花ID作为券码
	snowflakeID, err := generateSnowflakeID()
	if err != nil {
		tx.Rollback()
		return "", errors.New("生成券码失败: " + err.Error())
	}

	status := false // false 表示未使用
	now := time.Now()
	couponOrderUser := shop.CouponOrderUser{
		CouponNum:  snowflakeID, // 使用雪花ID作为券码
		CouponID:   &couponID,
		UserID:     userID,
		ShopUserID: &shopUserID,
		Status:     &status,
		ClaimedAt:  &now,
	}
	if err = tx.Create(&couponOrderUser).Error; err != nil {
		tx.Rollback()
		return "", errors.New("领取优惠券失败: " + err.Error())
	}

	// 4. 更新优惠券已领取数量 (使用 GORM 的表达式来增加 claimed 字段，并添加乐观锁或行锁避免并发问题)
	// result := tx.Model(&shop.Coupon{}).Where("id = ? AND claimed < quantity", couponID).UpdateColumn("claimed", gorm.Expr("claimed + ?", 1))
	// 更安全的做法是先锁定记录
	result := tx.Model(&coupon).Where("id = ?", couponID).Update("claimed", gorm.Expr("claimed + 1"))
	if result.Error != nil {
		tx.Rollback()
		return "", errors.New("更新优惠券数量失败: " + result.Error.Error())
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return "", errors.New("优惠券已领完或更新失败")
	}

	return snowflakeID, tx.Commit().Error
}

// IssueCouponToAllUsers 管理员向所有用户发放优惠券
// 注意：此函数可能需要根据实际用户表结构进行调整，这里假设用户表为 system.SysUser
func (couService *CouponOrderUserService) IssueCouponToAllUsers(ctx context.Context, couponID int) (err error) {
	// 1. 检查优惠券是否存在且有效
	var coupon shop.Coupon
	err = global.GVA_DB.Where("id = ? AND status = ? AND start_time <= NOW() AND end_time >= NOW()", couponID, true).First(&coupon).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("优惠券无效")
		}
		return errors.New("查询优惠券失败: " + err.Error())
	}

	// 2. 获取所有用户ID (假设用户表为 system_users, 字段为 id)
	// 您需要根据您的实际用户表结构修改此部分
	var userIDs []uint
	// 假设 shop_user 表存在并且有 id 字段
	if err = global.GVA_DB.Table("shop_user").Pluck("id", &userIDs).Error; err != nil {
		// 如果 shop_user 不存在，尝试 system.SysUser
		// 注意：gin-vue-admin 的用户表通常是 sys_users, 模型是 system.SysUser
		// 如果您的用户体系不同，请修改这里的查询
		if err = global.GVA_DB.Model(&system.SysUser{}).Pluck("id", &userIDs).Error; err != nil {
			return errors.New("获取用户列表失败: " + err.Error())
		}
	}

	if len(userIDs) == 0 {
		return errors.New("没有找到任何用户")
	}

	tx := global.GVA_DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err = tx.Error; err != nil {
		return err
	}

	issuedCount := 0
	status := false // false 表示未使用

	for _, userID := range userIDs {
		// 检查优惠券剩余数量
		if coupon.Quantity != nil && coupon.Claimed != nil && *coupon.Claimed+issuedCount >= *coupon.Quantity {
			global.GVA_LOG.Warn("优惠券库存不足，停止发放", zap.Int("couponID", couponID))
			break // 库存不足，停止发放
		}

		// 检查用户是否已领取
		var existingClaim shop.CouponOrderUser
		shopUserID := int(userID)
		findErr := tx.Where("coupon_id = ? AND shop_user_id = ?", couponID, shopUserID).First(&existingClaim).Error
		if findErr == nil {
			// 已领取，跳过
			continue
		} else if findErr != gorm.ErrRecordNotFound {
			tx.Rollback()
			return errors.New("检查用户领取状态失败: " + findErr.Error())
		}

		// 生成雪花ID作为券码
		snowflakeID, err := generateSnowflakeID()
		if err != nil {
			tx.Rollback()
			return errors.New("生成券码失败: " + err.Error())
		}

		couponOrderUser := shop.CouponOrderUser{
			CouponNum:  snowflakeID, // 使用雪花ID作为券码
			CouponID:   &couponID,
			ShopUserID: &shopUserID,
			Status:     &status,
		}
		if err = tx.Create(&couponOrderUser).Error; err != nil {
			tx.Rollback()
			return errors.New("为用户发放优惠券失败: " + err.Error())
		}
		issuedCount++
	}

	if issuedCount > 0 {
		// 更新优惠券已领取数量
		result := tx.Model(&shop.Coupon{}).Where("id = ?", couponID).UpdateColumn("claimed", gorm.Expr("claimed + ?", issuedCount))
		if result.Error != nil {
			tx.Rollback()
			return errors.New("批量更新优惠券数量失败: " + result.Error.Error())
		}
		if result.RowsAffected == 0 {
			// 理论上不应该发生，因为前面已经创建了 couponOrderUser 记录
			tx.Rollback()
			return errors.New("批量更新优惠券数量影响行数为0")
		}
	}

	return tx.Commit().Error
}

// couponAllowedProductIDs 解析优惠券允许的商品ID集合
// 优先使用 ProductIDs（逗号分隔），兼容旧字段 ProductID（单个）
func couponAllowedProductIDs(coupon shop.Coupon) map[int]bool {
	ids := make(map[int]bool)
	if coupon.ProductIDs != "" {
		parts := strings.Split(coupon.ProductIDs, ",")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if id, err := strconv.Atoi(p); err == nil && id > 0 {
				ids[id] = true
			}
		}
	}
	if len(ids) == 0 && coupon.ProductID != nil && *coupon.ProductID > 0 {
		ids[*coupon.ProductID] = true
	}
	return ids
}
