// 自动生成模板Address
package client

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 用户地址 结构体  Address
type Address struct {
	global.GVA_MODEL
	Name        string `json:"name" form:"name" gorm:"column:name;comment:收件人名称;" binding:"required"`                          //收件人名称
	Phone       string `json:"phone" form:"phone" gorm:"column:phone;comment:收件人电话;" binding:"required"`                       //收件人电话
	Province    *int   `json:"province" form:"province" gorm:"column:province;comment:收件省份;" binding:"required"`               //收件省份
	ProvinceStr string `json:"provinceStr" form:"provinceStr" gorm:"column:province_str;comment:收件省份（回显）;" binding:"required"` //收件省份
	City        *int   `json:"city" form:"city" gorm:"column:city;comment:收件城市;" binding:"required"`                           //收件城市
	CityStr     string `json:"cityStr" form:"cityStr" gorm:"column:city_str;comment:收件城市（回显）;" binding:"required"`             //收件城市
	Area        *int   `json:"area" form:"area" gorm:"column:area;comment:收件区;" binding:"required"`                            //收件区域
	AreaStr     string `json:"areaStr" form:"areaStr" gorm:"column:area_str;comment:收件区（回显）;" binding:"required"`              //收件区域
	Street      string `json:"street" form:"street" gorm:"column:street;comment:收件详细地址;" binding:"required"`                   //收件详细地址
	UserID      uint   `json:"userID" form:"userID" gorm:"column:user_id;comment:用户id;" `                                      //用户id
	Active      *bool  `json:"active" form:"active" gorm:"column:active;comment:是否默认;"`                                        //是否默认
}

// TableName 用户地址 Address自定义表名 client_address
func (Address) TableName() string {
	return "client_address"
}
