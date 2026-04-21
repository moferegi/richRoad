package initialize

import (
	"context"
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/model"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Gorm 自动迁移客服系统数据库表
func Gorm(ctx context.Context) {
	err := global.GVA_DB.WithContext(ctx).AutoMigrate(
		new(model.CsConversation),
		new(model.CsMessage),
		new(model.CsAgent),
		new(model.CsQuickReply),
		new(model.CsBlacklist),
		new(model.CsConfig),
	)
	if err != nil {
		err = errors.Wrap(err, "客服系统注册表失败!")
		zap.L().Error(fmt.Sprintf("%+v", err))
	}
}
