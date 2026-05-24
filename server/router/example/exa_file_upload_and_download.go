package example

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type FileUploadAndDownloadRouter struct{}

func (e *FileUploadAndDownloadRouter) InitFileUploadAndDownloadPublicRouter(Router *gin.RouterGroup) {
	fileUploadAndDownloadPublicRouter := Router.Group("fileUploadAndDownload")
	fileUploadAndDownloadTicketIssueRouter := Router.Group("fileUploadAndDownload").Use(middleware.JWTAuth())
	{
		fileUploadAndDownloadPublicRouter.POST("uploadByTicket", exaFileUploadAndDownloadApi.UploadFileByTicket)                  // 使用票据上传文件
		fileUploadAndDownloadTicketIssueRouter.POST("createScanUploadTicket", exaFileUploadAndDownloadApi.CreateScanUploadTicket) // 生成扫码上传票据
	}
}

func (e *FileUploadAndDownloadRouter) InitFileUploadAndDownloadRouter(Router *gin.RouterGroup) {
	fileUploadAndDownloadRouter := Router.Group("fileUploadAndDownload")
	{
		fileUploadAndDownloadRouter.POST("upload", exaFileUploadAndDownloadApi.UploadFile)                                 // 上传文件
		fileUploadAndDownloadRouter.POST("getFileList", exaFileUploadAndDownloadApi.GetFileList)                           // 获取上传文件列表
		fileUploadAndDownloadRouter.POST("deleteFile", exaFileUploadAndDownloadApi.DeleteFile)                             // 删除指定文件
		fileUploadAndDownloadRouter.POST("editFileName", exaFileUploadAndDownloadApi.EditFileName)                         // 编辑文件名或者备注
		fileUploadAndDownloadRouter.POST("breakpointContinue", exaFileUploadAndDownloadApi.BreakpointContinue)             // 断点续传
		fileUploadAndDownloadRouter.GET("findFile", exaFileUploadAndDownloadApi.FindFile)                                  // 查询当前文件成功的切片
		fileUploadAndDownloadRouter.POST("breakpointContinueFinish", exaFileUploadAndDownloadApi.BreakpointContinueFinish) // 切片传输完成
		fileUploadAndDownloadRouter.POST("removeChunk", exaFileUploadAndDownloadApi.RemoveChunk)                           // 删除切片
		fileUploadAndDownloadRouter.POST("importURL", exaFileUploadAndDownloadApi.ImportURL)                               // 导入URL
		fileUploadAndDownloadRouter.POST("signURL", exaFileUploadAndDownloadApi.SignURL)                                   // 生成防盗链签名URL
		fileUploadAndDownloadRouter.GET("hotlinkConfig", exaFileUploadAndDownloadApi.GetHotlinkConfig)                     // 获取防盗链配置信息
		fileUploadAndDownloadRouter.GET("listFolders", exaFileUploadAndDownloadApi.ListOSSFolders)                         // 列举OSS文件夹
	}
}
