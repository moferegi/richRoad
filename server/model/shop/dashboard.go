package shop

// DashboardOverview 数据看板概览
type DashboardOverview struct {
	// 销售
	TotalSales int64 `json:"totalSales"` // 总销售额(分)
	TodaySales int64 `json:"todaySales"` // 今日销售额

	// 订单
	OrderTotal     int64 `json:"orderTotal"`     // 总订单量
	OrderPending   int64 `json:"orderPending"`   // 待付款
	OrderPaid      int64 `json:"orderPaid"`      // 已付款(含已发/已收)
	OrderShipped   int64 `json:"orderShipped"`   // 已发货
	OrderReceived  int64 `json:"orderReceived"`  // 已收货
	OrderCancelled int64 `json:"orderCancelled"` // 已取消
	OrderRefunding int64 `json:"orderRefunding"` // 退款中
	OrderRefunded  int64 `json:"orderRefunded"`  // 已退款
	OrderToday     int64 `json:"orderToday"`     // 今日订单

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

	// 访客
	VisitorPV   int64 `json:"visitorPV"`   // 今日PV
	VisitorUV   int64 `json:"visitorUV"`   // 今日UV
	VisitorNew  int64 `json:"visitorNew"`  // 今日新访客
	SignInToday int64 `json:"signInToday"` // 今日签到人数

	// 商品
	ProductTotal  int64 `json:"productTotal"`  // 商品总数
	ProductActive int64 `json:"productActive"` // 上架商品

	// 趋势数据(近7天)
	SalesTrend []DayValue `json:"salesTrend"` // 每日销售额趋势
	OrderTrend []DayValue `json:"orderTrend"` // 每日订单数趋势
	UserTrend  []DayValue `json:"userTrend"`  // 每日新用户趋势
}

// DayValue 日期+数值对
type DayValue struct {
	Date  string `json:"date"`
	Value int64  `json:"value"`
}
