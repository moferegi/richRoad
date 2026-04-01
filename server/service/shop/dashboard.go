package shop

import (
	"strconv"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
)

type DashboardService struct{}

// GetDashboardOverview 获取数据看板概览
func (s *DashboardService) GetDashboardOverview() (overview shop.DashboardOverview, err error) {
	db := global.GVA_DB

	// 总销售额(已付款订单)
	db.Model(&shop.Order{}).Where("status IN ?", []string{"1", "2", "3"}).
		Select("COALESCE(SUM(total_price),0)").Scan(&overview.TotalSales)

	// 今日销售额
	today := time.Now().Format("2006-01-02")
	db.Model(&shop.Order{}).Where("status IN ? AND DATE(created_at) = ?", []string{"1", "2", "3"}, today).
		Select("COALESCE(SUM(total_price),0)").Scan(&overview.TodaySales)

	// 订单统计
	db.Model(&shop.Order{}).Count(&overview.OrderTotal)
	db.Model(&shop.Order{}).Where("status IN ?", []string{"1", "2", "3"}).Count(&overview.OrderPaid)
	db.Model(&shop.Order{}).Where("status = ?", "5").Count(&overview.OrderCancelled)

	// 优惠券统计
	db.Model(&shop.Coupon{}).Select("COALESCE(SUM(quantity),0)").Scan(&overview.CouponIssuedCount)
	db.Model(&shop.Coupon{}).Select("COALESCE(SUM(claimed),0)").Scan(&overview.CouponClaimedCount)
	db.Model(&shop.Coupon{}).Select("COALESCE(SUM(quantity * discount),0)").Scan(&overview.CouponIssuedAmount)

	// 已使用优惠券总金额
	var usedCouponAmount int64
	db.Model(&shop.CouponOrderUser{}).Where("status = ?", true).
		Joins("LEFT JOIN shop_coupon ON coupon_order_user.coupon_id = shop_coupon.id").
		Select("COALESCE(SUM(shop_coupon.discount),0)").Scan(&usedCouponAmount)
	overview.CouponUsedAmount = usedCouponAmount

	// 积分统计
	var pointsIssued, pointsUsed int64
	db.Model(&client.PointRecord{}).Where("change_type = ?", "increase").
		Select("COALESCE(SUM(point_change),0)").Scan(&pointsIssued)
	db.Model(&client.PointRecord{}).Where("change_type = ?", "decrease").
		Select("COALESCE(SUM(ABS(point_change)),0)").Scan(&pointsUsed)
	overview.PointsIssued = pointsIssued
	overview.PointsUsed = pointsUsed

	// 积分对应货币额(从配置读取汇率)
	var config client.SysConfig
	if e := db.Where("config_key = ?", "points_exchange_rate").First(&config).Error; e == nil {
		// 格式: "100:1" -> 100积分=1货币
		// 简单处理: 用已使用积分/比率
		rate := parseExchangeRate(config.ConfigValue)
		if rate > 0 {
			overview.PointsCurrencyValue = pointsUsed / int64(rate)
		}
	}

	// 用户统计
	db.Model(&client.ClientUser{}).Count(&overview.UserTotal)
	db.Model(&client.ClientUser{}).Where("DATE(created_at) = ?", today).Count(&overview.UserToday)

	return
}

func parseExchangeRate(value string) int {
	// 格式: "100:1" -> 返回100
	for i, c := range value {
		if c == ':' {
			rate, err := strconv.Atoi(value[:i])
			if err != nil {
				return 0
			}
			return rate
		}
	}
	rate, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return rate
}
