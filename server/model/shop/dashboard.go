package shop

// DashboardOverview 数据看板概览
type DashboardOverview struct {
	// 销售
	TotalSales int64 `json:"totalSales"` // 总销售额(分)
	TodaySales int64 `json:"todaySales"` // 今日销售额

	// 订单
	OrderTotal     int64 `json:"orderTotal"`     // 总订单量
	OrderPaid      int64 `json:"orderPaid"`      // 已付款
	OrderCancelled int64 `json:"orderCancelled"` // 已取消

	// 优惠券
	CouponIssuedCount  int64 `json:"couponIssuedCount"`  // 发放量(总数量)
	CouponClaimedCount int64 `json:"couponClaimedCount"` // 已领取量
	CouponIssuedAmount int64 `json:"couponIssuedAmount"` // 发放总金额
	CouponUsedAmount   int64 `json:"couponUsedAmount"`   // 已使用总金额

	// 积分
	PointsIssued        int64 `json:"pointsIssued"`        // 已发放积分
	PointsUsed          int64 `json:"pointsUsed"`          // 已使用积分
	PointsCurrencyValue int64 `json:"pointsCurrencyValue"` // 积分对应货币额

	// 用户
	UserTotal int64 `json:"userTotal"` // 总用户量
	UserToday int64 `json:"userToday"` // 今日新增
}
