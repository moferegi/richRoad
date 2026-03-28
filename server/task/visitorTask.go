package task

import (
	"fmt"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	"gorm.io/gorm"
)

const visitorLogRetainDays = 90

// VisitorDailyTask 每日执行：聚合昨日数据 + 清理旧日志
func VisitorDailyTask(db *gorm.DB) {
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	// 1. 聚合昨日汇总
	if err := aggregateSummary(db, yesterday); err != nil {
		fmt.Println("visitor aggregate error:", err)
	}

	// 2. 清理90天前的日志
	cutoff := time.Now().AddDate(0, 0, -visitorLogRetainDays).Format("2006-01-02")
	result := db.Where("created_date < ?", cutoff).Delete(&client.VisitorLog{})
	if result.Error != nil {
		fmt.Println("visitor cleanup error:", result.Error)
	} else if result.RowsAffected > 0 {
		fmt.Printf("visitor cleanup: deleted %d old logs (before %s)\n", result.RowsAffected, cutoff)
	}
}

func aggregateSummary(db *gorm.DB, date string) error {
	var pv int64
	db.Model(&client.VisitorLog{}).Where("created_date = ?", date).Count(&pv)

	var uv int64
	db.Model(&client.VisitorLog{}).Where("created_date = ?", date).Distinct("visitor_id").Count(&uv)

	var registeredUser int64
	db.Model(&client.VisitorLog{}).Where("created_date = ?", date).Where("user_id > 0").Distinct("user_id").Count(&registeredUser)

	var newVisitor int64
	db.Model(&client.VisitorLog{}).
		Where("created_date = ? AND visitor_id NOT IN (?)",
			date,
			db.Model(&client.VisitorLog{}).Select("visitor_id").Where("created_date < ?", date),
		).Distinct("visitor_id").Count(&newVisitor)

	summary := client.VisitorSummary{
		Date:           date,
		PV:             pv,
		UV:             uv,
		NewVisitor:     newVisitor,
		RegisteredUser: registeredUser,
	}

	return db.Where("date = ?", date).Assign(summary).FirstOrCreate(&summary).Error
}
