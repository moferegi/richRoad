package client

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	clientReq "github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
)

type AddressService struct {
}

// CreateAddress 创建用户地址记录
// Author [piexlmax](https://github.com/piexlmax)
func (addressService *AddressService) CreateAddress(address *client.Address) (err error) {
	err = global.GVA_DB.Create(address).Error
	return err
}

// DeleteAddress 删除用户地址记录
// Author [piexlmax](https://github.com/piexlmax)
func (addressService *AddressService) DeleteAddress(ID string, UserID uint) (err error) {
	err = global.GVA_DB.Delete(&client.Address{}, "id = ? and user_id = ?", ID, UserID).Error
	return err
}

// DeleteAddressByIds 批量删除用户地址记录
// Author [piexlmax](https://github.com/piexlmax)
func (addressService *AddressService) DeleteAddressByIds(IDs []string, UserID uint) (err error) {
	err = global.GVA_DB.Delete(&[]client.Address{}, "id in ? and user_id = ?", IDs, UserID).Error
	return err
}

// UpdateAddress 更新用户地址记录
// Author [piexlmax](https://github.com/piexlmax)
func (addressService *AddressService) UpdateAddress(address client.Address) (err error) {
	if *address.Active {
		// 如果用户设置默认地址，那么其他的都要改为非默认
		global.GVA_DB.Model(&client.Address{}).Where("user_id = ?", address.UserID).Update("active", false)
	}
	err = global.GVA_DB.Model(&client.Address{}).Where("id = ? and user_id = ?", address.ID, address.UserID).Updates(&address).Error
	return err
}

// GetAddress 根据ID获取用户地址记录
// Author [piexlmax](https://github.com/piexlmax)
func (addressService *AddressService) GetAddress(ID string, UserID uint) (address client.Address, err error) {
	err = global.GVA_DB.Where("id = ? and user_id = ?", ID, UserID).First(&address).Error
	return
}

// GetDefaultAddress 获取用户默认地址 如果用户没有设置则拉取表内最后一条
func (addressService *AddressService) GetDefaultAddress(UserID uint) (address client.Address, err error) {
	// 尝试获取默认地址
	err = global.GVA_DB.Where("user_id = ? and active = ?", UserID, true).First(&address).Error
	if err != nil {
		// 如果所有active都为0，则使用Order()按表内所有时间排序后返回最后一条
		err = global.GVA_DB.Where("user_id = ?", UserID).Order("created_at desc").First(&address).Error
	}
	return address, nil
}

// GetAddressInfoList 分页获取用户地址记录
// Author [piexlmax](https://github.com/piexlmax)
func (addressService *AddressService) GetAddressInfoList(info clientReq.AddressSearch) (list []client.Address, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&client.Address{})
	var addresss []client.Address
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.Name != "" {
		db = db.Where("name LIKE ?", "%"+info.Name+"%")
	}
	if info.UserID != 0 {
		db = db.Where("user_id = ?", info.UserID)
	}

	// 排序支持
	orderClause := "active desc, id desc"
	if info.OrderBy != "" {
		allowedCols := map[string]bool{"created_at": true}
		if allowedCols[info.OrderBy] {
			dir := "asc"
			if info.OrderDir == "desc" {
				dir = "desc"
			}
			orderClause = info.OrderBy + " " + dir
		}
	}
	db = db.Order(orderClause)

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&addresss).Error
	return addresss, total, err
}
func (addressService *AddressService) GetAddressDataSource() (res map[string][]map[string]any, err error) {
	res = make(map[string][]map[string]any)

	area := make([]map[string]any, 0)
	global.GVA_DB.Table("geo_areas").Select("name as label,code as value,name_i18n as labelI18n").Scan(&area)
	res["area"] = area
	city := make([]map[string]any, 0)
	global.GVA_DB.Table("geo_cities").Select("name as label,code as value,name_i18n as labelI18n").Scan(&city)
	res["city"] = city
	province := make([]map[string]any, 0)
	global.GVA_DB.Table("geo_provinces").Select("name as label,code as value,name_i18n as labelI18n").Scan(&province)
	res["province"] = province
	return
}
