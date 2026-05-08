package client

import (
	"errors"
	"sort"
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
		return errors.New("signInAlreadyToday")
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

// GetSignInList 管理端获取所有用户签到记录（支持分页、搜索、排序）
func (s *SignInService) GetSignInList(page, pageSize int, userId string, username string, startDate, endDate string, orderKey string, desc bool) (list []map[string]interface{}, total int64, err error) {
	// 子查询：每个用户的总签到天数
	totalDaysSub := "(SELECT COUNT(*) FROM client_sign_in AS t WHERE t.user_id = si.user_id)"
	// 主查询
	db := global.GVA_DB.Table("client_sign_in AS si").
		Select("si.id, si.user_id, si.sign_date, si.created_at, cu.username, cu.nickname, cu.phone, " + totalDaysSub + " AS total_days")

	db = db.Joins("LEFT JOIN client_user AS cu ON cu.id = si.user_id")

	if userId != "" {
		db = db.Where("si.user_id = ?", userId)
	}
	if username != "" {
		db = db.Where("cu.username LIKE ?", "%"+username+"%")
	}
	if startDate != "" {
		db = db.Where("si.sign_date >= ?", startDate)
	}
	if endDate != "" {
		db = db.Where("si.sign_date <= ?", endDate)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if pageSize > 0 {
		db = db.Limit(pageSize).Offset(pageSize * (page - 1))
	}

	// 排序
	orderStr := "si.sign_date DESC, si.id DESC"
	sortByContinuous := false
	if orderKey != "" {
		direction := "ASC"
		if desc {
			direction = "DESC"
		}
		switch orderKey {
		case "sign_date":
			orderStr = "si.sign_date " + direction + ", si.id DESC"
		case "total_days":
			orderStr = "total_days " + direction + ", si.sign_date DESC"
		case "continuous_days":
			sortByContinuous = true
		default:
			orderStr = "si.sign_date DESC, si.id DESC"
		}
	}

	err = db.Order(orderStr).Find(&list).Error
	if err != nil {
		return
	}

	// 为每条记录计算连续签到天数
	for i, item := range list {
		uid, ok := item["user_id"]
		if !ok {
			continue
		}
		var userID uint
		switch v := uid.(type) {
		case int64:
			userID = uint(v)
		case uint:
			userID = v
		case float64:
			userID = uint(v)
		}
		if userID > 0 {
			list[i]["continuous_days"] = s.GetContinuousSignInDays(userID)
		}
	}

	// 如果按连续签到天数排序，在内存中排序
	if sortByContinuous && len(list) > 0 {
		sort.Slice(list, func(i, j int) bool {
			ci, _ := list[i]["continuous_days"].(int)
			cj, _ := list[j]["continuous_days"].(int)
			if desc {
				return ci > cj
			}
			return ci < cj
		})
	}

	return
}

// DeleteSignIn 管理端删除签到记录
func (s *SignInService) DeleteSignIn(id uint) (err error) {
	err = global.GVA_DB.Delete(&client.SignIn{}, id).Error
	return
}
