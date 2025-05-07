// 自动生成模板Good
package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"gorm.io/datatypes"
)

// 商品 结构体  Good
type Good struct {
	global.GVA_MODEL
	Description string         `json:"description" form:"description" gorm:"column:description;comment:商品描述;"`        //商品描述
	Tags        datatypes.JSON `json:"tags" form:"tags" gorm:"column:tags;comment:商品标签;"`                             //商品标签
	Specs       datatypes.JSON `json:"specs" form:"specs" gorm:"column:specs;comment:商品规格;"`                          //商品规格
	Attrs       datatypes.JSON `json:"attrs" form:"attrs" gorm:"column:attrs;comment:商品属性;"`                          //商品属性
	ImageUrl    string         `json:"imageUrl" form:"imageUrl" gorm:"column:image_url;comment:商品图片URL;"`             //商品图片URL
	Banner      datatypes.JSON `json:"banner" form:"banner" gorm:"column:banner;comment:商品轮播图;"`                      //商品轮播图
	Price       *float64       `json:"price" form:"price" gorm:"column:price;comment:商品价格;"`                          //商品价格
	Rating      *float64       `json:"rating" form:"rating" gorm:"column:rating;comment:商品评分;"`                       //商品评分
	ReviewCount *int           `json:"reviewCount" form:"reviewCount" gorm:"column:review_count;comment:商品评论数量;"`     //商品评论数量
	SaleCount   *int           `json:"saleCount" form:"saleCount" gorm:"column:sale_count;comment:商品销售数量;"`           //商品销售数量
	Title       string         `json:"title" form:"title" gorm:"column:title;comment:商品名称;size:191;"`                 //商品名称
	CategoryID  *int           `json:"categoryID" form:"categoryID" gorm:"column:category_id;comment:商品类型;size:191;"` //商品类型
	Status      *bool          `json:"status" form:"status" gorm:"column:status;comment:状态;"`                         //状态
	Postage     *float64       `json:"postage" form:"postage" gorm:"column:postage;comment:邮费;"`                      //邮费
	Discount    *int           `json:"discount" form:"discount" gorm:"column:discount;comment:折扣;"`                   //折扣
	SKUS        []Sku          `json:"skus" gorm:"foreignKey:GoodID;references:ID"`                                   //SKU
	Detail      string         `json:"detail" form:"detail" gorm:"column:detail;comment:商品详情;type:text"`              //商品详情
	CollectNum  int            `json:"collect_num" form:"collect_num" gorm:"column:collect_num;comment:收藏数量;"`        //收藏数量
	SaleNum     uint           `json:"saleNum" form:"saleNum" gorm:"column:sale_num;comment:销量;"`                     //销量
}

// TableName 商品 Good自定义表名 shop_good
func (Good) TableName() string {
	return "shop_good"
}
