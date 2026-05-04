package initialize

import (
	"context"

	model "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/plugin-tool/utils"
)

// Api 注册客服系统 API 到 sys_apis 表
func Api(ctx context.Context) {
	entities := []model.SysApi{
		// 公共（WebSocket）
		{Path: "/cs/ws", Description: "客服WebSocket(用户)", ApiGroup: "客服系统", Method: "GET"},
		{Path: "/cs/wsAgent", Description: "客服WebSocket(坐席)", ApiGroup: "客服系统", Method: "GET"},
		// 客户端用户接口
		{Path: "/cs/conversation/getOrCreate", Description: "发起/获取会话", ApiGroup: "客服系统", Method: "POST"},
		{Path: "/cs/conversation/rate", Description: "评价会话", ApiGroup: "客服系统", Method: "POST"},
		{Path: "/cs/message/history", Description: "获取消息历史", ApiGroup: "客服系统", Method: "GET"},
		{Path: "/cs/message/upload", Description: "上传聊天图片", ApiGroup: "客服系统", Method: "POST"},
		{Path: "/cs/message/revoke", Description: "撤回消息", ApiGroup: "客服系统", Method: "POST"},
		{Path: "/cs/quickReply/all", Description: "获取全部快捷回复", ApiGroup: "客服系统", Method: "GET"},
		// 坐席接口
		{Path: "/cs/agent/conversation/list", Description: "坐席会话列表", ApiGroup: "客服系统", Method: "GET"},
		{Path: "/cs/agent/conversation/accept", Description: "坐席接入会话", ApiGroup: "客服系统", Method: "POST"},
		{Path: "/cs/agent/conversation/close", Description: "关闭会话", ApiGroup: "客服系统", Method: "POST"},
		{Path: "/cs/agent/conversation/transfer", Description: "转接会话", ApiGroup: "客服系统", Method: "POST"},
		{Path: "/cs/agent/message/send", Description: "坐席发送消息", ApiGroup: "客服系统", Method: "POST"},
		{Path: "/cs/agent/message/upload", Description: "坐席上传聊天图片", ApiGroup: "客服系统", Method: "POST"},
		{Path: "/cs/agent/quickReply/list", Description: "快捷回复列表(坐席)", ApiGroup: "客服系统", Method: "GET"},
		// 管理员接口
		{Path: "/cs/admin/conversation/list", Description: "所有会话列表", ApiGroup: "客服系统", Method: "GET"},
		{Path: "/cs/admin/conversation/assign", Description: "手动分配会话", ApiGroup: "客服系统", Method: "POST"},
		{Path: "/cs/admin/agent/list", Description: "坐席列表", ApiGroup: "客服系统", Method: "GET"},
		{Path: "/cs/admin/agent/create", Description: "创建坐席", ApiGroup: "客服系统", Method: "POST"},
		{Path: "/cs/admin/agent/update", Description: "更新坐席", ApiGroup: "客服系统", Method: "PUT"},
		{Path: "/cs/admin/agent/delete", Description: "删除坐席", ApiGroup: "客服系统", Method: "DELETE"},
		{Path: "/cs/admin/quickReply/create", Description: "创建快捷回复", ApiGroup: "客服系统", Method: "POST"},
		{Path: "/cs/admin/quickReply/update", Description: "更新快捷回复", ApiGroup: "客服系统", Method: "PUT"},
		{Path: "/cs/admin/quickReply/delete", Description: "删除快捷回复", ApiGroup: "客服系统", Method: "DELETE"},
		{Path: "/cs/admin/blacklist/list", Description: "黑名单列表", ApiGroup: "客服系统", Method: "GET"},
		{Path: "/cs/admin/blacklist/add", Description: "加入黑名单", ApiGroup: "客服系统", Method: "POST"},
		{Path: "/cs/admin/blacklist/remove", Description: "移出黑名单", ApiGroup: "客服系统", Method: "DELETE"},
		{Path: "/cs/config/get", Description: "获取客服配置(公开)", ApiGroup: "客服系统", Method: "GET"},
		{Path: "/cs/admin/config/get", Description: "获取客服配置(后台)", ApiGroup: "客服系统", Method: "GET"},
		{Path: "/cs/admin/config/update", Description: "更新客服配置", ApiGroup: "客服系统", Method: "PUT"},
	}
	utils.RegisterApis(entities...)
}
