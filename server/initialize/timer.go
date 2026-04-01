package initialize

import (
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/task"

	"github.com/robfig/cron/v3"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

func Timer() {
	go func() {
		var option []cron.Option
		option = append(option, cron.WithSeconds())
		// 清理DB定时任务
		_, err := global.GVA_Timer.AddTaskByFunc("ClearDB", "@daily", func() {
			err := task.ClearTable(global.GVA_DB) // 定时任务方法定在task文件包中
			if err != nil {
				fmt.Println("timer error:", err)
			}
		}, "定时清理数据库【日志，黑名单】内容", option...)
		if err != nil {
			fmt.Println("add timer error:", err)
		}

		// 访客日志：每天凌晨2点自动聚合昨日数据 + 清理90天前日志
		_, err = global.GVA_Timer.AddTaskByFunc("VisitorDailyTask", "0 0 2 * * *", func() {
			if global.GVA_DB == nil {
				return
			}
			task.VisitorDailyTask(global.GVA_DB)
		}, "访客统计日聚合+日志清理", option...)
		if err != nil {
			fmt.Println("add visitor timer error:", err)
		}

		// 其他定时任务定在这里 参考上方使用方法

		_, err = global.GVA_Timer.AddTaskByFunc("ClearOrder", "0 */5 * * * *", func() {
			if global.GVA_DB == nil {
				return
			}
			e := task.ClearOrder(global.GVA_DB)
			if e != nil {
				fmt.Println("订单清理失败:", e)
			}
			e = task.ClearExpiredOrders(global.GVA_DB)
			if e != nil {
				fmt.Println("过期订单清理失败:", e)
			}
		}, "定时清理已过期订单", option...)
		if err != nil {
			fmt.Println("add timer error:", err)
		}
	}()
}
