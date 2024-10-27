package client

import (
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	clientReq "github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"gorm.io/gorm"
)

type ClientUserService struct {
}

func (clientUserService *ClientUserService) Login(loginInfo *clientReq.Login) (clientUser client.ClientUser, err error) {
	err = global.GVA_DB.Where("username = ?", loginInfo.Username).First(&clientUser).Error
	if err != nil {
		return clientUser, errors.New("用户不存在")
	}
	if !utils.BcryptCheck(loginInfo.Password, clientUser.Password) {
		return clientUser, errors.New("密码错误")
	}
	return
}

// CreateClientUser 创建客户端用户记录
// Author [piexlmax](https://github.com/piexlmax)
func (clientUserService *ClientUserService) CreateClientUser(clientUser *client.ClientUser) (err error) {
	ferr := global.GVA_DB.Where("username = ?", clientUser.Username).First(&client.ClientUser{}).Error
	if ferr == nil {
		return errors.New("用户名已注册")
	}
	if clientUser.Nickname == "" {
		clientUser.Nickname = clientUser.Username
	}
	clientUser.Password = utils.BcryptHash(clientUser.Password)
	err = global.GVA_DB.Create(clientUser).Error
	return err
}

// DeleteClientUser 删除客户端用户记录
// Author [piexlmax](https://github.com/piexlmax)
func (clientUserService *ClientUserService) DeleteClientUser(ID string, userID uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&client.ClientUser{}).Where("id = ?", ID).Update("deleted_by", userID).Error; err != nil {
			return err
		}
		if err = tx.Delete(&client.ClientUser{}, "id = ?", ID).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteClientUserByIds 批量删除客户端用户记录
// Author [piexlmax](https://github.com/piexlmax)
func (clientUserService *ClientUserService) DeleteClientUserByIds(IDs []string, deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&client.ClientUser{}).Where("id in ?", IDs).Update("deleted_by", deleted_by).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", IDs).Delete(&client.ClientUser{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdateClientUser 更新客户端用户记录
// Author [piexlmax](https://github.com/piexlmax)
func (clientUserService *ClientUserService) UpdateClientUser(clientUser client.ClientUser) (err error) {
	var oldUser client.ClientUser
	err = global.GVA_DB.First(&oldUser, "id = ?", clientUser.ID).Updates(&clientUser).Error
	return err
}

// GetClientUser 根据ID获取客户端用户记录
// Author [piexlmax](https://github.com/piexlmax)
func (clientUserService *ClientUserService) GetClientUser(ID string) (clientUser client.ClientUser, err error) {
	err = global.GVA_DB.Omit("password").Where("id = ?", ID).First(&clientUser).Error
	return
}

// GetClientUserInfoList 分页获取客户端用户记录
// Author [piexlmax](https://github.com/piexlmax)
func (clientUserService *ClientUserService) GetClientUserInfoList(info clientReq.ClientUserSearch) (list []client.ClientUser, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&client.ClientUser{})
	var clientUsers []client.ClientUser
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.Username != "" {
		db = db.Where("username = ?", info.Username)
	}
	if info.Nickname != "" {
		db = db.Where("nickname LIKE ?", "%"+info.Nickname+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&clientUsers).Error
	return clientUsers, total, err
}

func (clientUserService *ClientUserService) SetClientUserInfo(key string, value string, userID uint) (err error) {
	err = global.GVA_DB.Model(&client.ClientUser{}).Where("id = ?", userID).Update(key, value).Error
	return err
}
