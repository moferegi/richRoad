package shop

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// GoodPurchase 进货统计 结构体
type GoodPurchase struct {
	global.GVA_MODEL
	GoodID       uint       `json:"goodID" form:"goodID" gorm:"column:good_id;index;comment:商品ID;"`                            //商品ID
	SkuID        *uint      `json:"skuID" form:"skuID" gorm:"column:sku_id;index;comment:SKU_ID;"`                              //SKU ID
	Quantity     int        `json:"quantity" form:"quantity" gorm:"column:quantity;comment:进货数量;"`                               //进货数量
	UnitCost     int        `json:"unitCost" form:"unitCost" gorm:"column:unit_cost;comment:单价(分);"`                            //单价(分)
	TotalCost    int        `json:"totalCost" form:"totalCost" gorm:"column:total_cost;comment:总金额(分);"`                        //总金额(分)
	Supplier     string     `json:"supplier" form:"supplier" gorm:"column:supplier;size:200;comment:供应商;"`                       //供应商
	PurchaseDate *time.Time `json:"purchaseDate" form:"purchaseDate" gorm:"column:purchase_date;comment:进货日期;"`                  //进货日期
	Remark       string     `json:"remark" form:"remark" gorm:"column:remark;size:500;comment:备注;"`                              //备注
	Good         Good       `json:"good" gorm:"foreignKey:GoodID;references:ID"`                                                //商品
	Sku          *Sku       `json:"sku" gorm:"foreignKey:SkuID;references:ID"`                                                  //SKU
}

// TableName 进货统计自定义表名
func (GoodPurchase) TableName() string {
	return "shop_good_purchase"
}
