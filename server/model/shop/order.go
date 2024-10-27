// 自动生成模板Order
package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"time"
)

// 订单 结构体  Order
type Order struct {
	global.GVA_MODEL
	UserID     uint          `json:"userID" form:"userID" gorm:"column:user_id;comment:购买者ID;"`               //购买者ID
	TotalPrice uint          `json:"totalPrice" form:"totalPrice" gorm:"column:total_price;comment:订单价格;"`    //订单价格（分）
	Express    string        `json:"express" form:"express" gorm:"column:express;comment:快递单号;"`              //快递单号
	Status     string        `json:"status" form:"status" gorm:"column:status;comment:订单状态;"`                 //订单状态
	Detail     []OrderDetail `json:"detail" gorm:"foreignKey:OrderID;references:ID"`                          //订单详情
	PayOrderID string        `json:"payOrderID" form:"payOrderID" gorm:"column:pay_order_id;comment:支付订单ID;"` //支付订单ID
	Phone      string        `json:"phone" form:"phone" gorm:"column:phone;comment:收件人电话;"`                   //收件人电话
	Name       string        `json:"name" form:"name" gorm:"column:name;comment:收件人姓名;"`                      //收件人姓名
	Province   string        `json:"province" form:"province" gorm:"column:province;comment:收件省份（回显）;"`       //收件省份
	City       string        `json:"city" form:"city" gorm:"column:city;comment:收件城市（回显）;"`                   //收件城市
	Area       string        `json:"area" form:"area" gorm:"column:area;comment:收件区（回显）;"`                    //收件区域
	Street     string        `json:"street" form:"street" gorm:"column:street;comment:收件详细地址;"`               //收件详细地址
	CloseTime  time.Time     `json:"closeTime" form:"closeTime" gorm:"column:close_time;comment:关闭时间;"`       //关闭时间
	Comment    *Comment      `json:"comment" gorm:"->;foreignKey:OrderID;references:ID"`
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
	UserID     uint             `json:"userID" form:"userID" gorm:"column:user_id;comment:购买者ID;"`            //购买者ID
	TotalPrice uint             `json:"totalPrice" form:"totalPrice" gorm:"column:total_price;comment:订单价格;"` //订单价格（分）
	Status     string           `json:"status" form:"status" gorm:"column:status;comment:订单状态;"`              //订单状态
	Detail     []OrderDetailRes `json:"detail" gorm:"foreignKey:OrderID;references:ID"`                       //订单详情
	Express    string           `json:"express" form:"express" gorm:"column:express;comment:快递单号;"`           //快递单号
	Phone      string           `json:"phone" form:"phone" gorm:"column:phone;comment:收件人电话;"`                //收件人电话
	Name       string           `json:"name" form:"name" gorm:"column:name;comment:收件人姓名;"`                   //收件人姓名
	Province   string           `json:"province" form:"province" gorm:"column:province;comment:收件省份（回显）;"`    //收件省份
	City       string           `json:"city" form:"city" gorm:"column:city;comment:收件城市（回显）;"`                //收件城市
	Area       string           `json:"area" form:"area" gorm:"column:area;comment:收件区（回显）;"`                 //收件区域
	Street     string           `json:"street" form:"street" gorm:"column:street;comment:收件详细地址;"`            //收件详细地址
	CloseTime  time.Time        `json:"closeTime" form:"closeTime" gorm:"column:close_time;comment:关闭时间;"`    //关闭时间
	Comment    *Comment         `json:"comment" gorm:"->;foreignKey:OrderID;references:ID"`
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
