// 自动生成模板Good
package shop

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"gorm.io/datatypes"
)

// 商品 结构体  Good
type Good struct {
	global.GVA_MODEL
	Description string         `json:"description" form:"description" gorm:"column:description;comment:商品描述;type:text"` //商品描述(JSON多语言)
	Tags        datatypes.JSON `json:"tags" form:"tags" gorm:"column:tags;comment:商品标签;"`                               //商品标签
	Specs       datatypes.JSON `json:"specs" form:"specs" gorm:"column:specs;comment:商品规格;"`                            //商品规格(JSON多语言key)
	Attrs       datatypes.JSON `json:"attrs" form:"attrs" gorm:"column:attrs;comment:商品属性;"`                            //商品属性(JSON多语言key)
	ImageUrl    string         `json:"imageUrl" form:"imageUrl" gorm:"column:image_url;comment:商品图片URL;"`               //商品图片URL
	UpperImage  string         `json:"upperImage" form:"upperImage" gorm:"column:upper_image;comment:上装图URL;size:500;"` //上装图URL
	LowerImage  string         `json:"lowerImage" form:"lowerImage" gorm:"column:lower_image;comment:下装图URL;size:500;"` //下装图URL
	Banner      datatypes.JSON `json:"banner" form:"banner" gorm:"column:banner;comment:商品轮播图;type:text"`               //商品轮播图
	Price       *float64       `json:"price" form:"price" gorm:"column:price;comment:商品价格;"`                            //商品价格
	PriceI18n   string         `json:"priceI18n" form:"priceI18n" gorm:"column:price_i18n;comment:多语言价格JSON;type:text"` //多语言价格(JSON)
	Rating      *float64       `json:"rating" form:"rating" gorm:"column:rating;comment:商品评分;"`                         //商品评分
	ReviewCount *int           `json:"reviewCount" form:"reviewCount" gorm:"column:review_count;comment:商品评论数量;"`       //商品评论数量
	SaleCount   *int           `json:"saleCount" form:"saleCount" gorm:"column:sale_count;comment:商品销售数量;"`             //商品销售数量
	Title       string         `json:"title" form:"title" gorm:"column:title;comment:商品名称;type:text"`                   //商品名称(JSON多语言)
	CategoryID  *int           `json:"categoryID" form:"categoryID" gorm:"column:category_id;comment:商品类型;size:191;"`   //商品类型
	Status      *bool          `json:"status" form:"status" gorm:"column:status;comment:状态;"`                           //状态
	Postage     *float64       `json:"postage" form:"postage" gorm:"column:postage;comment:邮费;"`                        //邮费
	Discount    *int           `json:"discount" form:"discount" gorm:"column:discount;comment:折扣;"`                     //折扣
	SKUS        []Sku          `json:"skus" gorm:"foreignKey:GoodID;references:ID"`                                     //SKU
	Detail      string         `json:"detail" form:"detail" gorm:"column:detail;comment:商品详情;type:text"`                //商品详情(JSON多语言富文本)
	CollectNum  int            `json:"collect_num" form:"collect_num" gorm:"column:collect_num;comment:收藏数量;"`          //收藏数量
	SaleNum     uint           `json:"saleNum" form:"saleNum" gorm:"column:sale_num;comment:销量;"`                       //销量
	Recommend   *bool          `json:"recommend" form:"recommend" gorm:"column:recommend;comment:是否推荐;"`                //是否推荐
	ViewNum     int            `json:"view_num" form:"view_num" gorm:"column:view_num;comment:浏览量;"`                    //浏览量

	// === 备注 ===
	Remark string `json:"remark" form:"remark" gorm:"column:remark;size:500;comment:备注(商品来源等);"` //备注

	// === 外部链接 ===
	ExternalImagePath string `json:"externalImagePath" form:"externalImagePath" gorm:"column:external_image_path;size:500;comment:外部图片路径(优先);"` //外部图片路径

	// === 积分抵扣设置 ===
	PointsEnabled  *bool `json:"pointsEnabled" form:"pointsEnabled" gorm:"column:points_enabled;default:false;comment:是否允许积分抵扣;"`  //是否允许积分抵扣
	PointsMaxUse   *int  `json:"pointsMaxUse" form:"pointsMaxUse" gorm:"column:points_max_use;default:0;comment:最多可用积分数;"`         //最多可用积分数
	PointsUseTimes *int  `json:"pointsUseTimes" form:"pointsUseTimes" gorm:"column:points_use_times;default:1;comment:同商品可用积分次数;"` //同商品可用积分次数

	// === 预售设置 ===
	IsPresale           *bool      `json:"isPresale" form:"isPresale" gorm:"column:is_presale;default:false;comment:是否预售;"`                                       //是否预售
	PresaleQty          *int       `json:"presaleQty" form:"presaleQty" gorm:"column:presale_qty;default:0;comment:预售数量;"`                                        //预售数量
	PresaleSold         *int       `json:"presaleSold" form:"presaleSold" gorm:"column:presale_sold;default:0;comment:已售数量;"`                                     //已售数量
	PresaleStart        *time.Time `json:"presaleStart" form:"presaleStart" gorm:"column:presale_start;comment:预售开始时间;"`                                          //预售开始时间
	PresaleEnd          *time.Time `json:"presaleEnd" form:"presaleEnd" gorm:"column:presale_end;comment:预售结束时间;"`                                                //预售结束时间
	PresaleEnabled      *bool      `json:"presaleEnabled" form:"presaleEnabled" gorm:"column:presale_enabled;default:true;comment:预售开关;"`                         //预售开关
	PresaleSort         *int       `json:"presaleSort" form:"presaleSort" gorm:"column:presale_sort;default:0;comment:预售排序;"`                                     //预售排序
	PresalePopupEnabled *bool      `json:"presalePopupEnabled" form:"presalePopupEnabled" gorm:"column:presale_popup_enabled;default:false;comment:预售弹窗开关;"`      //预售弹窗开关
	PresalePopupTitle   string     `json:"presalePopupTitle" form:"presalePopupTitle" gorm:"column:presale_popup_title;type:text;comment:预售弹窗标题(JSON多语言);"`       //预售弹窗标题
	PresalePopupContent string     `json:"presalePopupContent" form:"presalePopupContent" gorm:"column:presale_popup_content;type:text;comment:预售弹窗内容(JSON多语言);"` //预售弹窗内容

	// === 进货统计关联 ===
	Purchases []GoodPurchase `json:"purchases" gorm:"foreignKey:GoodID"`
}

// TableName 商品 Good自定义表名 shop_good
func (Good) TableName() string {
	return "shop_good"
}
