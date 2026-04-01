// 自动生成模板Order
package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"gorm.io/datatypes"
	"time"
)

type CreateOrder struct {
	global.GVA_MODEL
	CouponNum   string        `json:"couponNum" form:"couponNum" gorm:"column:coupon_num;comment:优惠券编号;"`     //优惠券编号
	UserID      uint          `json:"userID" form:"userID" gorm:"column:user_id;comment:购买者ID;"`              //购买者ID
	OriginPrice uint          `json:"originPrice" form:"originPrice" gorm:"column:origin_price;comment:原价;"`  //原价（分）
	TotalPrice  uint          `json:"totalPrice" form:"totalPrice" gorm:"column:total_price;comment:订单价格;"`   //订单价格（分）
	Discount    uint          `json:"discount" form:"discount" gorm:"column:discount;comment:优惠金额;"`          //优惠金额（分）
	UsePoints   bool          `json:"usePoints" form:"usePoints" gorm:"column:use_points;comment:是否使用积分;"`    //是否使用积分
	PointsUsed  uint          `json:"pointsUsed" form:"pointsUsed" gorm:"column:points_used;comment:使用的积分数;"` //使用的积分数
	Express     string        `json:"express" form:"express" gorm:"column:express;comment:快递单号;"`             //快递单号
	Status      string        `json:"status" form:"status" gorm:"column:status;comment:订单状态;"`                //订单状态
	RefundReason    string         `json:"refundReason" form:"refundReason" gorm:"column:refund_reason;comment:退款原因;"`           //退款原因
	RefundImages    datatypes.JSON `json:"refundImages" form:"refundImages" gorm:"column:refund_images;comment:退款图片;"`          //退款图片
	RefundAppliedAt *time.Time     `json:"refundAppliedAt" form:"refundAppliedAt" gorm:"column:refund_applied_at;comment:退款申请时间;"` //退款申请时间
	RefundRemark    string         `json:"refundRemark" form:"refundRemark" gorm:"column:refund_remark;comment:退款处理备注;"`         //退款处理备注
	RefundHandledAt *time.Time     `json:"refundHandledAt" form:"refundHandledAt" gorm:"column:refund_handled_at;comment:退款处理时间;"` //退款处理时间
	Detail      []OrderDetail `json:"detail" gorm:"foreignKey:OrderID;references:ID"`                         //订单详情
	OutTradeNo  string        `json:"outTradeNo" form:"outTradeNo" gorm:"column:out_trade_no;comment:商户订单号;"` //商户订单号
	Phone       string        `json:"phone" form:"phone" gorm:"column:phone;comment:收件人电话;"`                  //收件人电话
	Name        string        `json:"name" form:"name" gorm:"column:name;comment:收件人姓名;"`                     //收件人姓名
	Province    string        `json:"province" form:"province" gorm:"column:province;comment:收件省份（回显）;"`      //收件省份
	City        string        `json:"city" form:"city" gorm:"column:city;comment:收件城市（回显）;"`                  //收件城市
	Area        string        `json:"area" form:"area" gorm:"column:area;comment:收件区（回显）;"`                   //收件区域
	Street      string        `json:"street" form:"street" gorm:"column:street;comment:收件详细地址;"`              //收件详细地址
	CloseTime   time.Time     `json:"closeTime" form:"closeTime" gorm:"column:close_time;comment:关闭时间;"`      //关闭时间
	Comment     *Comment      `json:"comment" gorm:"->;foreignKey:OrderID;references:ID"`
	Pointed     bool          `json:"pointed" form:"pointed" gorm:"column:pointed;comment:是否已积分;"` //是否已积分
	IsPresale   bool          `json:"isPresale" form:"isPresale" gorm:"column:is_presale;default:false;comment:是否预售订单;"`   //是否预售订单
	PayMethod   string        `json:"payMethod" form:"payMethod" gorm:"column:pay_method;size:20;comment:支付方式(contact/qrcode);"` //支付方式
	ExpireAt    *time.Time    `json:"expireAt" form:"expireAt" gorm:"column:expire_at;comment:订单过期时间;"`                    //订单过期时间
	PaidAt      *time.Time    `json:"paidAt" form:"paidAt" gorm:"column:paid_at;comment:确认支付时间;"`                           //确认支付时间
	ReceivedAt  *time.Time    `json:"receivedAt" form:"receivedAt" gorm:"column:received_at;comment:确认收货时间;"`               //确认收货时间
	CancelledAt *time.Time    `json:"cancelledAt" form:"cancelledAt" gorm:"column:cancelled_at;comment:取消时间;"`                     //取消时间
}

// TableName 订单 Order自定义表名 shop_order
func (CreateOrder) TableName() string {
	return "shop_order"
}

// 订单 结构体  Order
type Order struct {
	global.GVA_MODEL
	CouponNum   string        `json:"couponNum" form:"couponNum" gorm:"column:coupon_num;comment:优惠券编号;"`     //优惠券编号
	UserID      uint          `json:"userID" form:"userID" gorm:"column:user_id;comment:购买者ID;"`              //购买者ID
	OriginPrice uint          `json:"originPrice" form:"originPrice" gorm:"column:origin_price;comment:原价;"`  //原价（分）
	TotalPrice  int           `json:"totalPrice" form:"totalPrice" gorm:"column:total_price;comment:订单价格;"`   //订单价格（分）
	Discount    uint          `json:"discount" form:"discount" gorm:"column:discount;comment:优惠金额;"`          //优惠金额（分）
	UsePoints   bool          `json:"usePoints" form:"usePoints" gorm:"column:use_points;comment:是否使用积分;"`    //是否使用积分
	PointsUsed  uint          `json:"pointsUsed" form:"pointsUsed" gorm:"column:points_used;comment:使用的积分数;"` //使用的积分数
	Express     string        `json:"express" form:"express" gorm:"column:express;comment:快递单号;"`             //快递单号
	Status      string        `json:"status" form:"status" gorm:"column:status;comment:订单状态;"`                //订单状态
	RefundReason    string         `json:"refundReason" form:"refundReason" gorm:"column:refund_reason;comment:退款原因;"`           //退款原因
	RefundImages    datatypes.JSON `json:"refundImages" form:"refundImages" gorm:"column:refund_images;comment:退款图片;"`          //退款图片
	RefundAppliedAt *time.Time     `json:"refundAppliedAt" form:"refundAppliedAt" gorm:"column:refund_applied_at;comment:退款申请时间;"` //退款申请时间
	RefundRemark    string         `json:"refundRemark" form:"refundRemark" gorm:"column:refund_remark;comment:退款处理备注;"`         //退款处理备注
	RefundHandledAt *time.Time     `json:"refundHandledAt" form:"refundHandledAt" gorm:"column:refund_handled_at;comment:退款处理时间;"` //退款处理时间
	Detail      []OrderDetail `json:"detail" gorm:"foreignKey:OrderID;references:ID"`                         //订单详情
	OutTradeNo  string        `json:"outTradeNo" form:"outTradeNo" gorm:"column:out_trade_no;comment:商户订单号;"` //商户订单号
	Phone       string        `json:"phone" form:"phone" gorm:"column:phone;comment:收件人电话;"`                  //收件人电话
	Name        string        `json:"name" form:"name" gorm:"column:name;comment:收件人姓名;"`                     //收件人姓名
	Province    string        `json:"province" form:"province" gorm:"column:province;comment:收件省份（回显）;"`      //收件省份
	City        string        `json:"city" form:"city" gorm:"column:city;comment:收件城市（回显）;"`                  //收件城市
	Area        string        `json:"area" form:"area" gorm:"column:area;comment:收件区（回显）;"`                   //收件区域
	Street      string        `json:"street" form:"street" gorm:"column:street;comment:收件详细地址;"`              //收件详细地址
	CloseTime   time.Time     `json:"closeTime" form:"closeTime" gorm:"column:close_time;comment:关闭时间;"`      //关闭时间
	Comment     *Comment      `json:"comment" gorm:"->;foreignKey:OrderID;references:ID"`
	Pointed     bool          `json:"pointed" form:"pointed" gorm:"column:pointed;comment:是否已积分;"` //是否已积分
	IsPresale   bool          `json:"isPresale" form:"isPresale" gorm:"column:is_presale;default:false;comment:是否预售订单;"`   //是否预售订单
	PayMethod   string        `json:"payMethod" form:"payMethod" gorm:"column:pay_method;size:20;comment:支付方式(contact/qrcode);"` //支付方式
	ExpireAt    *time.Time    `json:"expireAt" form:"expireAt" gorm:"column:expire_at;comment:订单过期时间;"`                    //订单过期时间
	PaidAt      *time.Time    `json:"paidAt" form:"paidAt" gorm:"column:paid_at;comment:确认支付时间;"`                           //确认支付时间
	ReceivedAt  *time.Time    `json:"receivedAt" form:"receivedAt" gorm:"column:received_at;comment:确认收货时间;"`               //确认收货时间
	CancelledAt *time.Time    `json:"cancelledAt" form:"cancelledAt" gorm:"column:cancelled_at;comment:取消时间;"`                     //取消时间
}

// TableName 订单 Order自定义表名 shop_order
func (Order) TableName() string {
	return "shop_order"
}

// 订单详情 结构体 OrderDetail
type OrderDetail struct {
	global.GVA_MODEL
	OrderID   uint `json:"orderID" form:"orderID" gorm:"column:order_id;comment:订单ID;"`       //订单ID
	GoodID    uint `json:"goodID" form:"goodID" gorm:"column:good_id;comment:商品ID;"`          //商品ID
	SKUID     uint `json:"skuID" form:"skuID" gorm:"column:sku_id;comment:skuID;"`            //skuID
	Quantity  uint `json:"quantity" form:"quantity" gorm:"column:quantity;comment:数量;"`       //数量
	Price     uint `json:"price" form:"price" gorm:"column:price;comment:商品单价;"`              //商品单价
	IsComment bool `json:"isComment" form:"isComment" gorm:"column:is_comment;comment:是否评论;"` //是否评论
}

// TableName 订单详情 OrderDetail自定义表名 shop_order_detail
func (OrderDetail) TableName() string {
	return "shop_order_detail"
}

type OrderRes struct {
	global.GVA_MODEL
	CouponNum   string           `json:"couponNum" form:"couponNum" gorm:"column:coupon_num;comment:优惠券编号;"`    //优惠券编号
	UserID      uint             `json:"userID" form:"userID" gorm:"column:user_id;comment:购买者ID;"`             //购买者ID
	OriginPrice uint             `json:"originPrice" form:"originPrice" gorm:"column:origin_price;comment:原价;"` //原价（分）
	Discount    uint             `json:"discount" form:"discount" gorm:"column:discount;comment:优惠金额;"`         //优惠金额（分）
	TotalPrice  uint             `json:"totalPrice" form:"totalPrice" gorm:"column:total_price;comment:订单价格;"`  //订单价格（分）
	Status      string           `json:"status" form:"status" gorm:"column:status;comment:订单状态;"`               //订单状态
	RefundReason    string         `json:"refundReason" form:"refundReason" gorm:"column:refund_reason;comment:退款原因;"`           //退款原因
	RefundImages    datatypes.JSON `json:"refundImages" form:"refundImages" gorm:"column:refund_images;comment:退款图片;"`          //退款图片
	RefundAppliedAt *time.Time     `json:"refundAppliedAt" form:"refundAppliedAt" gorm:"column:refund_applied_at;comment:退款申请时间;"` //退款申请时间
	RefundRemark    string         `json:"refundRemark" form:"refundRemark" gorm:"column:refund_remark;comment:退款处理备注;"`         //退款处理备注
	RefundHandledAt *time.Time     `json:"refundHandledAt" form:"refundHandledAt" gorm:"column:refund_handled_at;comment:退款处理时间;"` //退款处理时间
	Detail      []OrderDetailRes `json:"detail" gorm:"foreignKey:OrderID;references:ID"`                        //订单详情
	Express     string           `json:"express" form:"express" gorm:"column:express;comment:快递单号;"`            //快递单号
	Phone       string           `json:"phone" form:"phone" gorm:"column:phone;comment:收件人电话;"`                 //收件人电话
	Name        string           `json:"name" form:"name" gorm:"column:name;comment:收件人姓名;"`                    //收件人姓名
	Province    string           `json:"province" form:"province" gorm:"column:province;comment:收件省份（回显）;"`     //收件省份
	City        string           `json:"city" form:"city" gorm:"column:city;comment:收件城市（回显）;"`                 //收件城市
	Area        string           `json:"area" form:"area" gorm:"column:area;comment:收件区（回显）;"`                  //收件区域
	Street      string           `json:"street" form:"street" gorm:"column:street;comment:收件详细地址;"`             //收件详细地址
	CloseTime   time.Time        `json:"closeTime" form:"closeTime" gorm:"column:close_time;comment:关闭时间;"`     //关闭时间
	Comment     *Comment         `json:"comment" gorm:"->;foreignKey:OrderID;references:ID"`
	IsPresale   bool             `json:"isPresale" form:"isPresale" gorm:"column:is_presale;default:false;comment:是否预售订单;"`   //是否预售订单
	PayMethod   string           `json:"payMethod" form:"payMethod" gorm:"column:pay_method;size:20;comment:支付方式(contact/qrcode);"` //支付方式
	ExpireAt    *time.Time       `json:"expireAt" form:"expireAt" gorm:"column:expire_at;comment:订单过期时间;"`                    //订单过期时间
	PaidAt      *time.Time       `json:"paidAt" form:"paidAt" gorm:"column:paid_at;comment:确认支付时间;"`                           //确认支付时间
	ReceivedAt  *time.Time       `json:"receivedAt" form:"receivedAt" gorm:"column:received_at;comment:确认收货时间;"`               //确认收货时间
	CancelledAt *time.Time       `json:"cancelledAt" form:"cancelledAt" gorm:"column:cancelled_at;comment:取消时间;"`                     //取消时间
}

// TableName 订单 Order自定义表名 shop_order
func (OrderRes) TableName() string {
	return "shop_order"
}

type OrderCommentRes struct {
	global.GVA_MODEL
	UserID     uint           `json:"userID" form:"userID" gorm:"column:user_id;comment:购买者ID;"`            //购买者ID
	TotalPrice uint           `json:"totalPrice" form:"totalPrice" gorm:"column:total_price;comment:订单价格;"` //订单价格（分）
	Status     string         `json:"status" form:"status" gorm:"column:status;comment:订单状态;"`              //订单状态
	RefundReason    string         `json:"refundReason" form:"refundReason" gorm:"column:refund_reason;comment:退款原因;"`           //退款原因
	RefundImages    datatypes.JSON `json:"refundImages" form:"refundImages" gorm:"column:refund_images;comment:退款图片;"`          //退款图片
	RefundAppliedAt *time.Time     `json:"refundAppliedAt" form:"refundAppliedAt" gorm:"column:refund_applied_at;comment:退款申请时间;"` //退款申请时间
	RefundRemark    string         `json:"refundRemark" form:"refundRemark" gorm:"column:refund_remark;comment:退款处理备注;"`         //退款处理备注
	RefundHandledAt *time.Time     `json:"refundHandledAt" form:"refundHandledAt" gorm:"column:refund_handled_at;comment:退款处理时间;"` //退款处理时间
	Detail     OrderDetailRes `json:"detail" gorm:"foreignKey:OrderID;references:ID"`                       //订单详情
	Express    string         `json:"express" form:"express" gorm:"column:express;comment:快递单号;"`           //快递单号
	Phone      string         `json:"phone" form:"phone" gorm:"column:phone;comment:收件人电话;"`                //收件人电话
	Name       string         `json:"name" form:"name" gorm:"column:name;comment:收件人姓名;"`                   //收件人姓名
	Province   string         `json:"province" form:"province" gorm:"column:province;comment:收件省份（回显）;"`    //收件省份
	City       string         `json:"city" form:"city" gorm:"column:city;comment:收件城市（回显）;"`                //收件城市
	Area       string         `json:"area" form:"area" gorm:"column:area;comment:收件区（回显）;"`                 //收件区域
	Street     string         `json:"street" form:"street" gorm:"column:street;comment:收件详细地址;"`            //收件详细地址
	CloseTime  time.Time      `json:"closeTime" form:"closeTime" gorm:"column:close_time;comment:关闭时间;"`    //关闭时间
	Comment    Comment        `json:"comment" gorm:"->;foreignKey:OrderID;references:ID"`
}

// TableName 订单 Order自定义表名 shop_order
func (OrderCommentRes) TableName() string {
	return "shop_order"
}

type OrderDetailRes struct {
	global.GVA_MODEL
	OrderID   uint `json:"orderID" form:"orderID" gorm:"column:order_id;comment:订单ID;"` //订单ID
	GoodID    uint `json:"goodID" form:"goodID" gorm:"column:good_id;comment:商品ID;"`    //商品ID
	Good      Good `json:"good" gorm:"foreignKey:GoodID;references:ID"`
	SKUID     uint `json:"skuID" form:"skuID" gorm:"column:sku_id;comment:skuID;"` //skuID
	SKU       Sku  `json:"sku" gorm:"foreignKey:SKUID;references:ID"`
	Quantity  uint `json:"quantity" form:"quantity" gorm:"column:quantity;comment:数量;"`       //数量
	Price     uint `json:"price" form:"price" gorm:"column:price;comment:商品单价;"`              //商品单价
	IsComment bool `json:"isComment" form:"isComment" gorm:"column:is_comment;comment:是否评论;"` //是否评论
}

func (OrderDetailRes) TableName() string {
	return "shop_order_detail"
}
