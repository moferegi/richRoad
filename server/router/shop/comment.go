package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type CommentRouter struct{}

// InitCommentRouter 初始化 用户评论 路由信息
func (s *CommentRouter) InitCommentRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	commentRouter := Router.Group("comment").Use(middleware.OperationRecord())
	commentRouterWithoutRecord := Router.Group("comment")
	commentRouterWithoutAuth := PublicRouter.Group("comment")
	{
		commentRouter.POST("createComment", commentApi.CreateComment)             // 新建用户评论
		commentRouter.DELETE("deleteComment", commentApi.DeleteComment)           // 删除用户评论
		commentRouter.DELETE("deleteCommentByIds", commentApi.DeleteCommentByIds) // 批量删除用户评论
		commentRouter.PUT("updateComment", commentApi.UpdateComment)              // 更新用户评论
	}
	{
		commentRouterWithoutRecord.GET("getComment", commentApi.GetComment)         // 根据ID获取用户评论
		commentRouterWithoutRecord.GET("getCommentList", commentApi.GetCommentList) // 获取用户评论列表
	}
	{
		commentRouterWithoutAuth.GET("findComment", commentApi.FindComment)           // 根据商品ID获取公开评论
		commentRouterWithoutAuth.GET("getCommentPublic", commentApi.GetCommentPublic) // 获取用户评论列表
	}
}
