package client

import (
	"errors"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
)

type SignInService struct{}

// DoSignIn 用户签到
func (s *SignInService) DoSignIn(userID uint) (err error) {
	today := time.Now().Format("2006-01-02")
	todayTime, _ := time.Parse("2006-01-02", today)

	// 检查今日是否已签到
	var count int64
	global.GVA_DB.Model(&client.SignIn{}).Where("user_id = ? AND sign_date = ?", userID, today).Count(&count)
	if count > 0 {
		return errors.New("今日已签到")
	}

	signIn := client.SignIn{
		UserID:   userID,
		SignDate: todayTime,
	}
	err = global.GVA_DB.Create(&signIn).Error
	if err != nil {
		return
	}

	return
}

// GetSignInStatus 获取用户今日签到状态
func (s *SignInService) GetSignInStatus(userID uint) (signed bool) {
	today := time.Now().Format("2006-01-02")
	var count int64
	global.GVA_DB.Model(&client.SignIn{}).Where("user_id = ? AND sign_date = ?", userID, today).Count(&count)
	return count > 0
}

// GetSignInRecords 获取用户签到记录
func (s *SignInService) GetSignInRecords(userID uint, page, pageSize int) (list []client.SignIn, total int64, err error) {
	db := global.GVA_DB.Model(&client.SignIn{}).Where("user_id = ?", userID)
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if pageSize > 0 {
		db = db.Limit(pageSize).Offset(pageSize * (page - 1))
	}
	err = db.Order("sign_date DESC").Find(&list).Error
	return
}

// GetSignInCount 获取签到总天数
func (s *SignInService) GetSignInCount(userID uint) (count int64) {
	global.GVA_DB.Model(&client.SignIn{}).Where("user_id = ?", userID).Count(&count)
	return
}

// GetContinuousSignInDays 获取连续签到天数
func (s *SignInService) GetContinuousSignInDays(userID uint) (days int) {
	today := time.Now()
	for i := 0; ; i++ {
		checkDate := today.AddDate(0, 0, -i).Format("2006-01-02")
		var count int64
		global.GVA_DB.Model(&client.SignIn{}).Where("user_id = ? AND sign_date = ?", userID, checkDate).Count(&count)
		if count == 0 {
			break
		}
		days++
	}
	return
}
