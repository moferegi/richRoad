package shop

import (
	"strconv"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DashboardService struct{}

func scopePointAsset(db *gorm.DB) *gorm.DB {
	return db.Where("asset_type = ? OR asset_type IS NULL", client.AssetTypePoint)
}

// scanOrLog 执行查询并记录错误，不中断看板整体加载
func scanOrLog(label string, query *gorm.DB) {
	if query.Error != nil {
		global.GVA_LOG.Error("看板查询失败: "+label, zap.Error(query.Error))
	}
}

// GetDashboardOverview 获取数据看板概览
func (s *DashboardService) GetDashboardOverview() (overview shop.DashboardOverview, err error) {
	db := global.GVA_DB
	today := time.Now().Format("2006-01-02")

	paidStatuses := []string{"1", "2", "3", "7"}

	// === 销售额 ===
	scanOrLog("totalSales", db.Model(&shop.Order{}).Where("status IN ?", paidStatuses).
		Select("COALESCE(SUM(total_price),0)").Scan(&overview.TotalSales))
	scanOrLog("todaySales", db.Model(&shop.Order{}).Where("status IN ? AND DATE(created_at) = ?", paidStatuses, today).
		Select("COALESCE(SUM(total_price),0)").Scan(&overview.TodaySales))

	// === 订单统计 ===
	scanOrLog("orderTotal", db.Model(&shop.Order{}).Count(&overview.OrderTotal))
	scanOrLog("orderPending", db.Model(&shop.Order{}).Where("status = ?", "0").Count(&overview.OrderPending))
	scanOrLog("orderPaid", db.Model(&shop.Order{}).Where("status IN ?", paidStatuses).Count(&overview.OrderPaid))
	scanOrLog("orderShipped", db.Model(&shop.Order{}).Where("status = ?", "2").Count(&overview.OrderShipped))
	scanOrLog("orderReceived", db.Model(&shop.Order{}).Where("status = ?", "3").Count(&overview.OrderReceived))
	scanOrLog("orderCancelled", db.Model(&shop.Order{}).Where("status = ?", "4").Count(&overview.OrderCancelled))
	scanOrLog("orderRefunding", db.Model(&shop.Order{}).Where("status = ?", "6").Count(&overview.OrderRefunding))
	scanOrLog("orderRefunded", db.Model(&shop.Order{}).Where("status = ?", "5").Count(&overview.OrderRefunded))
	scanOrLog("orderToday", db.Model(&shop.Order{}).Where("DATE(created_at) = ?", today).Count(&overview.OrderToday))

	// === 优惠券统计 ===
	scanOrLog("couponIssuedCount", db.Model(&shop.Coupon{}).Select("COALESCE(SUM(quantity),0)").Scan(&overview.CouponIssuedCount))
	scanOrLog("couponClaimedCount", db.Model(&shop.Coupon{}).Select("COALESCE(SUM(claimed),0)").Scan(&overview.CouponClaimedCount))
	scanOrLog("couponIssuedAmount", db.Model(&shop.Coupon{}).Select("COALESCE(SUM(quantity * discount),0)").Scan(&overview.CouponIssuedAmount))
	var usedCouponAmount int64
	scanOrLog("couponUsedAmount", db.Model(&shop.CouponOrderUser{}).Where("coupon_order_user.status = ?", true).
		Joins("LEFT JOIN shop_coupon ON coupon_order_user.coupon_id = shop_coupon.id").
		Joins("LEFT JOIN shop_order ON coupon_order_user.order_id = shop_order.id").
		Where("shop_order.status IN ?", paidStatuses).
		Select("COALESCE(SUM(shop_coupon.discount),0)").Scan(&usedCouponAmount))
	overview.CouponUsedAmount = usedCouponAmount

	// === 积分统计 ===
	var pointsIssued int64
	// 发放积分：所有增加记录，排除退还积分(point_refund)
	scanOrLog("pointsIssued", scopePointAsset(db.Model(&client.PointRecord{})).Where("change_type = ? AND operation_type != ?", "increase", "point_refund").
		Select("COALESCE(SUM(point_change),0)").Scan(&pointsIssued))
	// 消耗积分（净值）：point_exchange 总额 - point_refund 退还总额
	var pointExchangeTotal, pointRefundTotal int64
	scanOrLog("pointExchange", scopePointAsset(db.Model(&client.PointRecord{})).Where("operation_type = ?", "point_exchange").
		Select("COALESCE(SUM(ABS(point_change)),0)").Scan(&pointExchangeTotal))
	scanOrLog("pointRefund", scopePointAsset(db.Model(&client.PointRecord{})).Where("operation_type = ?", "point_refund").
		Select("COALESCE(SUM(ABS(point_change)),0)").Scan(&pointRefundTotal))
	pointsUsed := pointExchangeTotal - pointRefundTotal
	if pointsUsed < 0 {
		pointsUsed = 0
	}
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
	scanOrLog("userTotal", db.Model(&client.ClientUser{}).Count(&overview.UserTotal))
	scanOrLog("userToday", db.Model(&client.ClientUser{}).Where("DATE(created_at) = ?", today).Count(&overview.UserToday))

	// === 访客 & 签到 ===
	scanOrLog("visitorPV", db.Model(&client.VisitorLog{}).Where("created_date = ?", today).Count(&overview.VisitorPV))
	scanOrLog("visitorUV", db.Model(&client.VisitorLog{}).Where("created_date = ?", today).Distinct("visitor_id").Count(&overview.VisitorUV))
	scanOrLog("visitorNew", db.Model(&client.VisitorLog{}).
		Where("created_date = ? AND visitor_id NOT IN (?)",
			today,
			db.Model(&client.VisitorLog{}).Select("visitor_id").Where("created_date < ?", today),
		).Distinct("visitor_id").Count(&overview.VisitorNew))
	scanOrLog("signInToday", db.Model(&client.SignIn{}).Where("DATE(sign_date) = ?", today).Count(&overview.SignInToday))

	// === 商品统计 ===
	scanOrLog("productTotal", db.Model(&shop.Good{}).Count(&overview.ProductTotal))
	scanOrLog("productActive", db.Model(&shop.Good{}).Where("status = ?", true).Count(&overview.ProductActive))

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
			Where("status IN ? AND DATE(created_at) = ?", []string{"1", "2", "3", "7"}, d).
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
