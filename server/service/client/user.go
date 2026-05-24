package client

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	clientReq "github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClientUserService struct {
}

var clientUserPhoneLocks sync.Map

func lockClientUserPhone(areaCode, phone string) func() {
	key := strings.TrimSpace(areaCode) + ":" + strings.TrimSpace(phone)
	lockValue, _ := clientUserPhoneLocks.LoadOrStore(key, &sync.Mutex{})
	mu := lockValue.(*sync.Mutex)
	mu.Lock()
	return mu.Unlock
}

func (clientUserService *ClientUserService) Login(loginInfo *clientReq.Login) (clientUser client.ClientUser, err error) {
	err = global.GVA_DB.Where("username = ?", loginInfo.Username).First(&clientUser).Error
	if err != nil {
		return clientUser, errors.New("loginFail")
	}
	if clientUser.Banned != nil && *clientUser.Banned {
		return clientUser, errors.New("accountBanned")
	}
	if !utils.BcryptCheck(loginInfo.Password, clientUser.Password) {
		return clientUser, errors.New("loginFail")
	}
	return
}

// CreateClientUser 创建客户端用户记录
// Author [piexlmax](https://github.com/piexlmax)
func (clientUserService *ClientUserService) CreateClientUser(clientUser *client.ClientUser) (err error) {
	ferr := global.GVA_DB.Where("username = ?", clientUser.Username).First(&client.ClientUser{}).Error
	if ferr == nil {
		return errors.New("usernameOrPhoneExists")
	}
	if clientUser.Nickname == "" {
		clientUser.Nickname = clientUser.Username
	}
	clientUser.UUID, _ = uuid.NewUUID()
	clientUser.Password = utils.BcryptHash(clientUser.Password)
	clientUser.InviteCode = generateInviteCode()
	clientUser.Point = 0
	clientUser.TryonPoint = 0
	clientUser.OpenID = ""
	clientUser.InvitedBy = 0
	clientUser.SubOrderRewarded = false
	err = global.GVA_DB.Create(clientUser).Error
	return err
}

// generateInviteCode 生成8位随机邀请码
func generateInviteCode() string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
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
	updates := map[string]interface{}{
		"username":   clientUser.Username,
		"avatar":     clientUser.Avatar,
		"nickname":   clientUser.Nickname,
		"gender":     clientUser.Gender,
		"phone":      clientUser.Phone,
		"email":      clientUser.Email,
		"area_code":  clientUser.AreaCode,
		"banned":     clientUser.Banned,
		"updated_by": clientUser.UpdatedBy,
	}
	err = global.GVA_DB.Model(&client.ClientUser{}).Where("id = ?", clientUser.ID).Updates(updates).Error
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
	db := global.GVA_DB.Model(&client.ClientUser{}).Omit("password")
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

	// 排序支持
	orderClause := "id desc"
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

func (clientUserService *ClientUserService) SetClientUserPhone(userID uint, phone string) error {
	phone = strings.TrimSpace(phone)
	if userID == 0 || phone == "" {
		return errors.New("phoneRequired")
	}

	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var user client.ClientUser
		if err := tx.Where("id = ?", userID).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("userNotExist")
			}
			return err
		}

		areaCode := strings.TrimSpace(user.AreaCode)
		var phoneAreaCode client.PhoneAreaCode
		query := tx.Where("is_enabled = ?", true)
		if areaCode != "" {
			query = query.Where("area_code = ?", areaCode)
		}
		if err := query.Order("sort ASC, id ASC").First(&phoneAreaCode).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("areaCodeInvalid")
			}
			return err
		}
		areaCode = strings.TrimSpace(phoneAreaCode.AreaCode)
		if areaCode == "" {
			return errors.New("areaCodeInvalid")
		}

		unlock := lockClientUserPhone(areaCode, phone)
		defer unlock()

		if phoneAreaCode.PhoneRegex != "" {
			matched, err := regexp.MatchString(phoneAreaCode.PhoneRegex, phone)
			if err != nil || !matched {
				return errors.New("phoneFormatError")
			}
		}

		var count int64
		if err := tx.Model(&client.ClientUser{}).
			Where("area_code = ? AND phone = ? AND id <> ?", areaCode, phone, userID).
			Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("phoneAlreadyRegistered")
		}

		return tx.Model(&client.ClientUser{}).Where("id = ?", userID).Updates(map[string]interface{}{
			"area_code": areaCode,
			"phone":     phone,
		}).Error
	})
}

// AdjustTryonPoint 后台调整用户试衣币，统一写入试衣币流水
func (clientUserService *ClientUserService) AdjustTryonPoint(ctx context.Context, req clientReq.AdjustTryonPointRequest, operatorID uint) error {
	if req.UserID == 0 {
		return errors.New("userIDRequired")
	}
	if req.Amount <= 0 {
		return errors.New("adjustAmountMustPositive")
	}
	changeType := strings.TrimSpace(req.ChangeType)
	if changeType != "increase" && changeType != "decrease" {
		return errors.New("pointChangeTypeInvalid")
	}
	reason := strings.TrimSpace(req.Reason)
	if reason == "" {
		return errors.New("adjustReasonRequired")
	}

	uid := int(req.UserID)
	assetType := client.AssetTypeTryonPoint
	operationType := "admin_adjust_tryon_point"
	pointChange := req.Amount
	remark := strings.TrimSpace(req.Remark)
	if remark != "" {
		remark = fmt.Sprintf("operator:%d; %s", operatorID, remark)
	} else {
		remark = fmt.Sprintf("operator:%d", operatorID)
	}

	record := client.PointRecord{
		AssetType:     &assetType,
		UserId:        &uid,
		ChangeType:    &changeType,
		PointChange:   &pointChange,
		OperationType: &operationType,
		Reason:        &reason,
		Remark:        &remark,
	}

	pointRecordService := PointRecordService{}
	return pointRecordService.CreatePointRecord(ctx, &record)
}

// GetSubordinates 获取用户的直接下级列表
func (clientUserService *ClientUserService) GetSubordinates(userID uint, page, pageSize int) (list []client.ClientUser, total int64, err error) {
	db := global.GVA_DB.Model(&client.ClientUser{}).Where("invited_by = ?", userID)
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	offset := pageSize * (page - 1)
	err = db.Select("id, username, nickname, avatar, created_at, invite_code").
		Order("created_at desc").Limit(pageSize).Offset(offset).Find(&list).Error
	return
}

// GetSubordinateCount 获取下级总数（含间接下级）
func (clientUserService *ClientUserService) GetSubordinateCount(userID uint) (int64, error) {
	var count int64
	err := global.GVA_DB.Model(&client.ClientUser{}).Where("invited_by = ?", userID).Count(&count).Error
	return count, err
}

// LoginByPhone 手机号+密码登录
func (clientUserService *ClientUserService) LoginByPhone(areaCode, phone, password string) (clientUser client.ClientUser, err error) {
	err = global.GVA_DB.Where("area_code = ? AND phone = ?", areaCode, phone).First(&clientUser).Error
	if err != nil {
		return clientUser, errors.New("phoneLoginFail")
	}
	if clientUser.Banned != nil && *clientUser.Banned {
		return clientUser, errors.New("accountBanned")
	}
	if !utils.BcryptCheck(password, clientUser.Password) {
		return clientUser, errors.New("phoneLoginFail")
	}
	return
}

// RegisterByPhone 手机号注册
func (clientUserService *ClientUserService) RegisterByPhone(areaCode, phone, password string, inviteCode string) (clientUser client.ClientUser, err error) {
	areaCode = strings.TrimSpace(areaCode)
	phone = strings.TrimSpace(phone)
	inviteCode = strings.TrimSpace(inviteCode)
	if areaCode == "" || phone == "" || password == "" {
		return clientUser, errors.New("invalidParams")
	}

	unlock := lockClientUserPhone(areaCode, phone)
	defer unlock()

	// 检查手机号是否已注册
	var count int64
	global.GVA_DB.Model(&client.ClientUser{}).Where("area_code = ? AND phone = ?", areaCode, phone).Count(&count)
	if count > 0 {
		return clientUser, errors.New("phoneAlreadyRegistered")
	}

	// 生成随机用户名(8位)
	username := generateRandomUsername()
	// 确保用户名不重复
	for i := 0; i < 10; i++ {
		var c int64
		global.GVA_DB.Model(&client.ClientUser{}).Where("username = ?", username).Count(&c)
		if c == 0 {
			break
		}
		username = generateRandomUsername()
	}

	clientUser.UUID, _ = uuid.NewUUID()
	clientUser.Username = username
	clientUser.Nickname = maskedPhoneNickname(phone)
	clientUser.Password = utils.BcryptHash(password)
	clientUser.Phone = phone
	clientUser.AreaCode = areaCode
	clientUser.Avatar = "https://qmplusimg.henrongyi.top/gva_header.jpg"
	clientUser.InviteCode = generateInviteCode()

	// 处理邀请码
	if inviteCode != "" {
		var inviter client.ClientUser
		if err := global.GVA_DB.Where("invite_code = ?", inviteCode).First(&inviter).Error; err == nil {
			clientUser.InvitedBy = inviter.ID
		}
	}

	err = global.GVA_DB.Create(&clientUser).Error
	return
}

func maskedPhoneNickname(phone string) string {
	if len(phone) <= 5 {
		return phone
	}
	return phone[:3] + "****" + phone[len(phone)-2:]
}

// generateRandomUsername 生成随机8位数字用户名
func generateRandomUsername() string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	// 生成8位数字
	num := int(b[0])<<24 | int(b[1])<<16 | int(b[2])<<8 | int(b[3])
	if num < 0 {
		num = -num
	}
	return fmt.Sprintf("%08d", num%100000000)
}

// ChangePassword 修改密码
func (clientUserService *ClientUserService) ChangePassword(userID uint, oldPassword, newPassword, method string) error {
	if strings.TrimSpace(method) != "old_password" || strings.TrimSpace(oldPassword) == "" || strings.TrimSpace(newPassword) == "" {
		return errors.New("invalidParams")
	}

	var user client.ClientUser
	err := global.GVA_DB.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return errors.New("userNotExist")
	}

	if !utils.BcryptCheck(oldPassword, user.Password) {
		return errors.New("passwordError")
	}

	newHash := utils.BcryptHash(newPassword)
	return global.GVA_DB.Model(&client.ClientUser{}).Where("id = ?", userID).Update("password", newHash).Error
}
