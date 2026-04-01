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
	today := time.Now().Format("2006-01-02")

	// === 销售额 ===
	db.Model(&shop.Order{}).Where("status IN ?", []string{"1", "2", "3"}).
		Select("COALESCE(SUM(total_price),0)").Scan(&overview.TotalSales)
	db.Model(&shop.Order{}).Where("status IN ? AND DATE(created_at) = ?", []string{"1", "2", "3"}, today).
		Select("COALESCE(SUM(total_price),0)").Scan(&overview.TodaySales)

	// === 订单统计 ===
	db.Model(&shop.Order{}).Count(&overview.OrderTotal)
	db.Model(&shop.Order{}).Where("status = ?", "0").Count(&overview.OrderPending)
	db.Model(&shop.Order{}).Where("status IN ?", []string{"1", "2", "3"}).Count(&overview.OrderPaid)
	db.Model(&shop.Order{}).Where("status = ?", "2").Count(&overview.OrderShipped)
	db.Model(&shop.Order{}).Where("status = ?", "3").Count(&overview.OrderReceived)
	db.Model(&shop.Order{}).Where("status = ?", "4").Count(&overview.OrderCancelled)
	db.Model(&shop.Order{}).Where("status = ?", "6").Count(&overview.OrderRefunding)
	db.Model(&shop.Order{}).Where("status = ?", "5").Count(&overview.OrderRefunded)
	db.Model(&shop.Order{}).Where("DATE(created_at) = ?", today).Count(&overview.OrderToday)

	// === 优惠券统计 ===
	db.Model(&shop.Coupon{}).Select("COALESCE(SUM(quantity),0)").Scan(&overview.CouponIssuedCount)
	db.Model(&shop.Coupon{}).Select("COALESCE(SUM(claimed),0)").Scan(&overview.CouponClaimedCount)
	db.Model(&shop.Coupon{}).Select("COALESCE(SUM(quantity * discount),0)").Scan(&overview.CouponIssuedAmount)
	var usedCouponAmount int64
	db.Model(&shop.CouponOrderUser{}).Where("status = ?", true).
		Joins("LEFT JOIN shop_coupon ON coupon_order_user.coupon_id = shop_coupon.id").
		Select("COALESCE(SUM(shop_coupon.discount),0)").Scan(&usedCouponAmount)
	overview.CouponUsedAmount = usedCouponAmount

	// === 积分统计 ===
	var pointsIssued, pointsUsed int64
	db.Model(&client.PointRecord{}).Where("change_type = ?", "increase").
		Select("COALESCE(SUM(point_change),0)").Scan(&pointsIssued)
	db.Model(&client.PointRecord{}).Where("change_type = ?", "decrease").
		Select("COALESCE(SUM(ABS(point_change)),0)").Scan(&pointsUsed)
	overview.PointsIssued = pointsIssued
	overview.PointsUsed = pointsUsed

	var config client.SysConfig
	if e := db.Where("config_key = ?", "points_exchange_rate").First(&config).Error; e == nil {
		rate := parseExchangeRate(config.ConfigValue)
		if rate > 0 {
			overview.PointsCurrencyValue = pointsUsed / int64(rate)
		}
	}

	// === 用户统计 ===
	db.Model(&client.ClientUser{}).Count(&overview.UserTotal)
	db.Model(&client.ClientUser{}).Where("DATE(created_at) = ?", today).Count(&overview.UserToday)

	// === 访客 & 签到 ===
	db.Model(&client.VisitorLog{}).Where("created_date = ?", today).Count(&overview.VisitorPV)
	db.Model(&client.VisitorLog{}).Where("created_date = ?", today).Distinct("visitor_id").Count(&overview.VisitorUV)
	db.Model(&client.VisitorLog{}).
		Where("created_date = ? AND visitor_id NOT IN (?)",
			today,
			db.Model(&client.VisitorLog{}).Select("visitor_id").Where("created_date < ?", today),
		).Distinct("visitor_id").Count(&overview.VisitorNew)
	db.Model(&client.SignIn{}).Where("DATE(sign_date) = ?", today).Count(&overview.SignInToday)

	// === 商品统计 ===
	db.Model(&shop.Good{}).Count(&overview.ProductTotal)
	db.Model(&shop.Good{}).Where("status = ?", true).Count(&overview.ProductActive)

	// === 7天趋势 ===
	overview.SalesTrend = s.get7DaySalesTrend()
	overview.OrderTrend = s.get7DayOrderTrend()
	overview.UserTrend = s.get7DayUserTrend()

	return
}

func (s *DashboardService) get7DaySalesTrend() []shop.DayValue {
	db := global.GVA_DB
	result := make([]shop.DayValue, 7)
	for i := 6; i >= 0; i-- {
		d := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		var val int64
		db.Model(&shop.Order{}).
			Where("status IN ? AND DATE(created_at) = ?", []string{"1", "2", "3"}, d).
			Select("COALESCE(SUM(total_price),0)").Scan(&val)
		result[6-i] = shop.DayValue{Date: d, Value: val}
	}
	return result
}

func (s *DashboardService) get7DayOrderTrend() []shop.DayValue {
	db := global.GVA_DB
	result := make([]shop.DayValue, 7)
	for i := 6; i >= 0; i-- {
		d := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		var val int64
		db.Model(&shop.Order{}).Where("DATE(created_at) = ?", d).Count(&val)
		result[6-i] = shop.DayValue{Date: d, Value: val}
	}
	return result
}

func (s *DashboardService) get7DayUserTrend() []shop.DayValue {
	db := global.GVA_DB
	result := make([]shop.DayValue, 7)
	for i := 6; i >= 0; i-- {
		d := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		var val int64
		db.Model(&client.ClientUser{}).Where("DATE(created_at) = ?", d).Count(&val)
		result[6-i] = shop.DayValue{Date: d, Value: val}
	}
	return result
}

func parseExchangeRate(value string) int {
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
