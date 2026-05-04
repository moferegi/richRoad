package customer_service

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/initialize"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/service"
	interfaces "github.com/flipped-aurora/gin-vue-admin/server/utils/plugin/v2"
	"github.com/gin-gonic/gin"
)

var _ interfaces.Plugin = (*plugin)(nil)

// Plugin 客服系统插件入口
var Plugin = new(plugin)

type plugin struct{}

func init() {
	interfaces.Register(Plugin)
}

func (p *plugin) Register(group *gin.Engine) {
	ctx := context.Background()
	// 注册 API 元信息
	initialize.Api(ctx)
	// 注册侧边栏菜单
	initialize.Menu(ctx)
	// 自动迁移数据库表
	initialize.Gorm(ctx)
	// 注册 Casbin 权限规则（幂等）
	initialize.InitCasbin()
	// 挂载路由
	initialize.Router(group)
	// 启动空闲会话超时自动关闭后台任务
	service.Service.ConfigService.StartAutoCloseWorker()
	// 启动排队会话重试分配后台任务，降低偶发卡队列概率
	service.Service.ConversationService.StartPendingDispatchWorker()
}
