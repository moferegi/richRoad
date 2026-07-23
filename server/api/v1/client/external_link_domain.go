package client

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/i18n"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ExternalLinkDomainApi struct{}

// CreateExternalLinkDomain 创建外部链接域名
// @Tags ExternalLinkDomain
// @Summary 创建外部链接域名
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body client.ExternalLinkDomain true "域名信息"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /extDomain/createExternalLinkDomain [post]
func (a *ExternalLinkDomainApi) CreateExternalLinkDomain(c *gin.Context) {
	var domain client.ExternalLinkDomain
	if err := c.ShouldBindJSON(&domain); err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	if err := extDomainService.CreateExternalLinkDomain(domain); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "createFail"), c)
		return
	}
	response.OkWithMessage(i18n.T(c, "createSuccess"), c)
}

// DeleteExternalLinkDomain 删除外部链接域名
// @Tags ExternalLinkDomain
// @Summary 删除外部链接域名
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "ID"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /extDomain/deleteExternalLinkDomain [delete]
func (a *ExternalLinkDomainApi) DeleteExternalLinkDomain(c *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	if err := extDomainService.DeleteExternalLinkDomain(req.ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "deleteFail"), c)
		return
	}
	response.OkWithMessage(i18n.T(c, "deleteSuccess"), c)
}

// UpdateExternalLinkDomain 更新外部链接域名
// @Tags ExternalLinkDomain
// @Summary 更新外部链接域名
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body client.ExternalLinkDomain true "域名信息"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /extDomain/updateExternalLinkDomain [put]
func (a *ExternalLinkDomainApi) UpdateExternalLinkDomain(c *gin.Context) {
	var domain client.ExternalLinkDomain
	if err := c.ShouldBindJSON(&domain); err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	if err := extDomainService.UpdateExternalLinkDomain(domain); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "updateFail"), c)
		return
	}
	response.OkWithMessage(i18n.T(c, "updateSuccess"), c)
}

// GetExternalLinkDomainList 分页获取外部链接域名列表
// @Tags ExternalLinkDomain
// @Summary 分页获取外部链接域名列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.ExternalLinkDomainSearch true "分页参数"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /extDomain/getExternalLinkDomainList [get]
func (a *ExternalLinkDomainApi) GetExternalLinkDomainList(c *gin.Context) {
	var pageInfo request.ExternalLinkDomainSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	list, total, err := extDomainService.GetExternalLinkDomainList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "getFail"), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, i18n.T(c, "getSuccess"), c)
}

// SetDefaultDomain 设置默认域名
// @Tags ExternalLinkDomain
// @Summary 设置默认域名
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body object true "ID"
// @Success 200 {object} response.Response{msg=string} "设置成功"
// @Router /extDomain/setDefaultDomain [post]
func (a *ExternalLinkDomainApi) SetDefaultDomain(c *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	if err := extDomainService.SetDefaultDomain(req.ID); err != nil {
		global.GVA_LOG.Error("设置失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "setFail"), c)
		return
	}
	response.OkWithMessage(i18n.T(c, "setSuccess"), c)
}

// GetDefaultDomain 获取默认域名（公开接口）
// @Tags ExternalLinkDomain
// @Summary 获取默认域名
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=string,msg=string} "获取成功"
// @Router /extDomain/getDefaultDomain [get]
func (a *ExternalLinkDomainApi) GetDefaultDomain(c *gin.Context) {
	domain, err := extDomainService.GetDefaultDomain()
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "getFail"), c)
		return
	}
	response.OkWithDetailed(domain, i18n.T(c, "getSuccess"), c)
}

// PingCloud 检测云存储连接
// @Tags ExternalLinkDomain
// @Summary 检测云存储连接
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.CloudPingReq true "域名ID"
// @Success 200 {object} response.Response{data=request.CloudPingResp,msg=string} "检测结果"
// @Router /extDomain/pingCloud [post]
func (a *ExternalLinkDomainApi) PingCloud(c *gin.Context) {
	var req request.CloudPingReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	result, err := cloudStorageService.PingCloud(req)
	if err != nil {
		global.GVA_LOG.Error("检测连接失败!", zap.Error(err))
		response.FailWithMessage("检测失败: "+err.Error(), c)
		return
	}
	response.OkWithData(result, c)
}

// ListCloudFiles 列出云存储文件
// @Tags ExternalLinkDomain
// @Summary 列出云存储文件
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.CloudListFilesReq true "查询参数"
// @Success 200 {object} response.Response{data=request.CloudListFilesResp,msg=string} "文件列表"
// @Router /extDomain/listCloudFiles [get]
func (a *ExternalLinkDomainApi) ListCloudFiles(c *gin.Context) {
	var req request.CloudListFilesReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	result, err := cloudStorageService.ListCloudFiles(req)
	if err != nil {
		global.GVA_LOG.Error("列出文件失败!", zap.Error(err))
		response.FailWithMessage("列出文件失败: "+err.Error(), c)
		return
	}
	response.OkWithData(result, c)
}

// CompareDirectories 多云目录比对
// @Tags ExternalLinkDomain
// @Summary 多云目录比对
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.CloudCompareReq true "查询参数"
// @Success 200 {object} response.Response{data=request.CloudCompareResp,msg=string} "比对结果"
// @Router /extDomain/compareDirectories [get]
func (a *ExternalLinkDomainApi) CompareDirectories(c *gin.Context) {
	var req request.CloudCompareReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	result, err := cloudStorageService.CompareDirectories(req)
	if err != nil {
		global.GVA_LOG.Error("目录比对失败!", zap.Error(err))
		response.FailWithMessage("比对失败: "+err.Error(), c)
		return
	}
	response.OkWithData(result, c)
}

// DeleteCloudFiles 批量删除云存储文件
// @Tags ExternalLinkDomain
// @Summary 批量删除云存储文件
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.CloudDeleteFilesReq true "删除参数（含日期密码）"
// @Success 200 {object} response.Response{data=request.CloudDeleteFilesResp,msg=string} "删除结果"
// @Router /extDomain/deleteCloudFiles [post]
func (a *ExternalLinkDomainApi) DeleteCloudFiles(c *gin.Context) {
	var req request.CloudDeleteFilesReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	result, err := cloudStorageService.DeleteCloudFiles(req)
	if err != nil {
		global.GVA_LOG.Error("删除文件失败!", zap.Error(err))
		response.FailWithMessage("删除失败: "+err.Error(), c)
		return
	}
	response.OkWithData(result, c)
}

// UploadCloudFile 上传文件到指定云存储
// @Tags ExternalLinkDomain
// @Summary 上传文件到指定云存储
// @Security ApiKeyAuth
// @accept multipart/form-data
// @Produce application/json
// @Param id formData int true "域名ID"
// @Param file formData file true "文件"
// @Param folder formData string false "上传目录"
// @Success 200 {object} response.Response{data=object,msg=string} "上传成功"
// @Router /extDomain/uploadCloudFile [post]
func (a *ExternalLinkDomainApi) UploadCloudFile(c *gin.Context) {
	idStr := c.DefaultPostForm("id", "0")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		response.FailWithMessage("无效的域名ID", c)
		return
	}

	folder := strings.TrimSpace(c.DefaultPostForm("folder", ""))

	_, header, err := c.Request.FormFile("file")
	if err != nil {
		global.GVA_LOG.Error("接收文件失败!", zap.Error(err))
		response.FailWithMessage("接收文件失败", c)
		return
	}

	filePath, _, err := cloudStorageService.UploadFileToCloudByID(uint(id), header, folder)
	if err != nil {
		global.GVA_LOG.Error("上传文件失败!", zap.Error(err))
		response.FailWithMessage("上传失败: "+err.Error(), c)
		return
	}

	response.OkWithData(map[string]string{"url": filePath, "key": filePath}, c)
}

// SearchCloudFiles 全局搜索云存储文件
// @Tags ExternalLinkDomain
// @Summary 全局搜索云存储文件
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.SearchCloudFilesReq true "搜索参数"
// @Success 200 {object} response.Response{data=request.SearchCloudFilesResp,msg=string} "搜索结果"
// @Router /extDomain/searchCloudFiles [get]
func (a *ExternalLinkDomainApi) SearchCloudFiles(c *gin.Context) {
	var req request.SearchCloudFilesReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	result, err := cloudStorageService.SearchCloudFiles(req)
	if err != nil {
		global.GVA_LOG.Error("搜索文件失败!", zap.Error(err))
		response.FailWithMessage("搜索失败: "+err.Error(), c)
		return
	}
	response.OkWithData(result, c)
}

// GetFileDownloadURL 获取文件下载链接
// @Tags ExternalLinkDomain
// @Summary 获取文件下载链接
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.DownloadCloudFileReq true "下载参数"
// @Success 200 {object} response.Response{data=request.DownloadCloudFileResp,msg=string} "下载链接"
// @Router /extDomain/getFileDownloadURL [get]
func (a *ExternalLinkDomainApi) GetFileDownloadURL(c *gin.Context) {
	var req request.DownloadCloudFileReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	result, err := cloudStorageService.GetFileDownloadURL(req)
	if err != nil {
		global.GVA_LOG.Error("获取下载链接失败!", zap.Error(err))
		response.FailWithMessage("获取下载链接失败: "+err.Error(), c)
		return
	}
	response.OkWithData(result, c)
}

// DownloadCloudFolder 下载目录为 zip
// @Tags ExternalLinkDomain
// @Summary 下载目录为 zip
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/octet-stream
// @Param id query int true "域名ID"
// @Param prefix query string true "目录前缀"
// @Success 200 {file} binary "zip文件"
// @Router /extDomain/downloadCloudFolder [get]
func (a *ExternalLinkDomainApi) DownloadCloudFolder(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		response.FailWithMessage("无效的域名ID", c)
		return
	}
	prefix := strings.TrimSpace(c.Query("prefix"))

	reader, filename, err := cloudStorageService.ZipCloudFolderByID(uint(id), prefix)
	if err != nil {
		global.GVA_LOG.Error("打包目录失败!", zap.Error(err))
		response.FailWithMessage("打包失败: "+err.Error(), c)
		return
	}
	defer reader.Close()

	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Header("Content-Type", "application/zip")
	c.Status(200)

	_, copyErr := io.Copy(c.Writer, reader)
	if copyErr != nil {
		global.GVA_LOG.Error("传输zip文件失败", zap.Error(copyErr))
	}
}
