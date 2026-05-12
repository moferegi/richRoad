package task

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"gorm.io/gorm"
)

// RefreshProcessingTryonTasks 定时补偿刷新processing试衣任务状态
func RefreshProcessingTryonTasks(db *gorm.DB) error {
	if db == nil {
		return errors.New("db Cannot be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 55*time.Second)
	defer cancel()

	scanned, changed, err := service.ServiceGroupApp.ClientServiceGroup.TryonTaskService.RefreshProcessingTryonTasks(ctx, db, 120)
	if err != nil {
		return err
	}

	if scanned > 0 {
		fmt.Printf("[RefreshProcessingTryonTasks] 扫描=%d, 状态变更=%d\n", scanned, changed)
	}

	return nil
}
