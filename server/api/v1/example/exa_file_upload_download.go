package example

import (
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/example"
	"github.com/flipped-aurora/gin-vue-admin/server/model/example/request"
	exampleRes "github.com/flipped-aurora/gin-vue-admin/server/model/example/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/i18n"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/upload"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type FileUploadAndDownloadApi struct{}

// UploadFile
// @Tags      ExaFileUploadAndDownload
// @Summary   上传文件示例
// @Security  ApiKeyAuth
// @accept    multipart/form-data
// @Produce   application/json
// @Param     file  formData  file                                                           true  "上传文件示例"
// @Param     folder  formData  string                                                        false "上传目录，如 tryon/clothes/source"
// @Param     uploadType  formData  string                                                    false "上传类型标签，如 shoe/person/cloth"
// @Param     uploadPosition  formData  string                                                false "上传位置标签，如 tryon/kefu"
// @Success   200   {object}  response.Response{data=exampleRes.ExaFileResponse,msg=string}  "上传文件示例,返回包括文件详情"
// @Router    /fileUploadAndDownload/upload [post]
func (b *FileUploadAndDownloadApi) UploadFile(c *gin.Context) {
	var file example.ExaFileUploadAndDownload
	noSave := c.DefaultQuery("noSave", "0")
	operatorUserID := utils.GetUserID(c)
	folder := normalizeUploadFolder(c.DefaultPostForm("folder", c.DefaultQuery("folder", "")))
	uploadType := strings.TrimSpace(c.DefaultPostForm("uploadType", c.DefaultQuery("uploadType", "")))
	uploadPosition := strings.TrimSpace(c.DefaultPostForm("uploadPosition", c.DefaultQuery("uploadPosition", "")))
	_, header, err := c.Request.FormFile("file")
	classId, _ := strconv.Atoi(c.DefaultPostForm("classId", "0"))
	if err != nil {
		global.GVA_LOG.Error("接收文件失败!", zap.Error(err))
		response.FailWithMessage("接收文件失败", c)
		return
	}
	file, err = fileUploadAndDownloadService.UploadFile(header, noSave, classId, folder, uploadType, uploadPosition, operatorUserID) // 文件上传后拿到文件路径
	if err != nil {
		global.GVA_LOG.Error("上传文件失败!", zap.Error(err))
		response.FailWithMessage("上传文件失败", c)
		return
	}
	response.OkWithDetailed(exampleRes.ExaFileResponse{File: file}, "上传成功", c)
}

// CreateScanUploadTicket
// @Tags      ExaFileUploadAndDownload
// @Summary   生成扫码上传一次性票据
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data body request.CreateScanUploadTicketRequest true "分类/目录"
// @Success   200  {object} response.Response{data=exampleRes.ScanUploadTicketResponse,msg=string}  "生成成功"
// @Router    /fileUploadAndDownload/createScanUploadTicket [post]
func (b *FileUploadAndDownloadApi) CreateScanUploadTicket(c *gin.Context) {
	if !ensureFileLibraryManager(c) {
		return
	}

	operatorUserID := utils.GetUserID(c)
	if operatorUserID == 0 {
		response.FailWithMessage(i18n.T(c, "noPermission"), c)
		return
	}

	var req request.CreateScanUploadTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}

	folder := normalizeUploadFolder(req.Folder)
	uploadType := strings.TrimSpace(req.UploadType)
	uploadPosition := strings.TrimSpace(req.UploadPosition)
	ticket, expireAt, err := fileUploadAndDownloadService.CreateScanUploadTicket(req.ClassId, folder, uploadType, uploadPosition, operatorUserID)
	if err != nil {
		global.GVA_LOG.Error("生成扫码上传票据失败!", zap.Error(err))
		response.FailWithMessage("生成扫码上传票据失败", c)
		return
	}

	response.OkWithDetailed(exampleRes.ScanUploadTicketResponse{
		Ticket:    ticket,
		ExpiresAt: expireAt.UnixMilli(),
	}, "生成成功", c)
}

// UploadFileByTicket
// @Tags      ExaFileUploadAndDownload
// @Summary   使用扫码上传票据上传文件
// @accept    multipart/form-data
// @Produce   application/json
// @Param     file  formData  file   true  "上传文件示例"
// @Param     ticket  formData  string  true  "扫码上传一次性票据"
// @Success   200   {object}  response.Response{data=exampleRes.ExaFileResponse,msg=string}  "上传成功"
// @Router    /fileUploadAndDownload/uploadByTicket [post]
func (b *FileUploadAndDownloadApi) UploadFileByTicket(c *gin.Context) {
	var file example.ExaFileUploadAndDownload
	allowed, waitSeconds := fileUploadAndDownloadService.CheckScanUploadTicketUploadRateLimit(c.ClientIP(), c.GetHeader("User-Agent"))
	if !allowed {
		exampleBannedIPService.RecordAttack(c.ClientIP(), "scan_upload_ticket_rate_limit", "", "扫码上传票据上传频率超限")
		msg := i18n.T(c, "requestTooFrequent")
		if waitSeconds > 0 {
			msg = msg + " " + strconv.Itoa(waitSeconds) + "s"
		}
		response.FailWithMessage(msg, c)
		return
	}

	ticket := strings.TrimSpace(c.PostForm("ticket"))
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		global.GVA_LOG.Error("接收文件失败!", zap.Error(err))
		response.FailWithMessage("接收文件失败", c)
		return
	}

	payload, err := fileUploadAndDownloadService.ConsumeScanUploadTicket(ticket)
	if err != nil {
		exampleBannedIPService.RecordAttack(c.ClientIP(), "scan_upload_ticket_invalid", "", "扫码上传票据无效或已过期")
		response.FailWithMessage("上传凭证无效或已过期", c)
		return
	}

	file, err = fileUploadAndDownloadService.UploadFile(
		header,
		"0",
		payload.ClassId,
		normalizeUploadFolder(payload.Folder),
		strings.TrimSpace(payload.UploadType),
		strings.TrimSpace(payload.UploadPosition),
		payload.OperatorUserID,
	)
	if err != nil {
		global.GVA_LOG.Error("扫码上传文件失败!", zap.Error(err))
		response.FailWithMessage("上传文件失败", c)
		return
	}

	response.OkWithDetailed(exampleRes.ExaFileResponse{File: file}, "上传成功", c)
}

func normalizeUploadFolder(raw string) string {
	raw = strings.TrimSpace(strings.ReplaceAll(raw, "\\\\", "/"))
	if raw == "" {
		return ""
	}

	parts := strings.Split(raw, "/")
	cleaned := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" || part == "." || part == ".." {
			continue
		}
		cleaned = append(cleaned, part)
	}

	return strings.Join(cleaned, "/")
}

func isFileLibraryManagerRole(authorityID uint) bool {
	return authorityID == 888 || authorityID == 8881 || authorityID == 9528
}

var publicSignMediaExtensions = map[string]struct{}{
	".jpg":  {},
	".jpeg": {},
	".png":  {},
	".gif":  {},
	".webp": {},
	".mp4":  {},
	".mov":  {},
	".webm": {},
	".m3u8": {},
	".ts":   {},
}

func isPublicSignMediaPathAllowed(filePath string, userID uint) bool {
	normalizedPath := "/" + strings.TrimPrefix(strings.TrimSpace(filePath), "/")
	trimmedPath := strings.TrimPrefix(normalizedPath, "/")
	if len(trimmedPath) < 8 {
		return false
	}
	if _, ok := publicSignMediaExtensions[strings.ToLower(path.Ext(normalizedPath))]; !ok {
		return false
	}

	patterns := []string{"%" + normalizedPath + "%", "%" + trimmedPath + "%"}
	publicGoodDB := func() *gorm.DB {
		return global.GVA_DB.Table("shop_good AS g").
			Where("g.status = ?", true).
			Where("g.category_id IS NULL OR g.category_id NOT IN (SELECT id FROM shop_category WHERE show_in_uni = ?)", false)
	}

	previewQuery := `(g.image_url LIKE ? OR g.upper_image LIKE ? OR g.lower_image LIKE ? OR g.external_image_path LIKE ? OR g.banner LIKE ? OR EXISTS (
		SELECT 1 FROM shop_sku AS s
		WHERE s.good_id = g.id AND (s.picture LIKE ? OR s.upper_image LIKE ? OR s.lower_image LIKE ? OR s.external_picture_path LIKE ?)
	))`
	for _, pattern := range patterns {
		var total int64
		if err := publicGoodDB().Where(previewQuery,
			pattern, pattern, pattern, pattern, pattern,
			pattern, pattern, pattern, pattern,
		).Limit(1).Count(&total).Error; err == nil && total > 0 {
			return true
		}
	}

	type mediaGood struct {
		ID    uint
		Price *float64
	}
	protectedQuery := `(g.detail LIKE ? OR g.attrs LIKE ? OR g.specs LIKE ? OR EXISTS (
		SELECT 1 FROM shop_sku AS s
		WHERE s.good_id = g.id AND (s.description LIKE ? OR s.attrs LIKE ? OR s.specs LIKE ?)
	))`
	for _, pattern := range patterns {
		var goods []mediaGood
		if err := publicGoodDB().Select("g.id, g.price").Where(protectedQuery,
			pattern, pattern, pattern, pattern, pattern, pattern,
		).Limit(20).Scan(&goods).Error; err != nil {
			continue
		}
		for _, good := range goods {
			if good.ID == 0 {
				continue
			}
			if good.Price == nil || *good.Price <= 0 || userHasPaidGood(userID, good.ID) {
				return true
			}
		}
	}
	return false
}

func userHasPaidGood(userID uint, goodID uint) bool {
	if userID == 0 || goodID == 0 {
		return false
	}
	var total int64
	err := global.GVA_DB.Table("shop_order AS o").
		Joins("JOIN shop_order_detail AS d ON d.order_id = o.id").
		Where("o.user_id = ? AND d.good_id = ?", userID, goodID).
		Where("o.paid_at IS NOT NULL AND o.status IN ?", []string{"1", "2", "3", "7"}).
		Limit(1).
		Count(&total).Error
	return err == nil && total > 0
}

func ensureFileLibraryManager(c *gin.Context) bool {
	if !isFileLibraryManagerRole(utils.GetUserAuthorityId(c)) {
		response.FailWithMessage(i18n.T(c, "noPermission"), c)
		return false
	}
	return true
}

func normalizeSignFilePath(raw string) (string, bool) {
	raw = strings.TrimSpace(strings.ReplaceAll(raw, "\\", "/"))
	if raw == "" || len(raw) > 2048 {
		return "", false
	}

	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		parsed, err := url.Parse(raw)
		if err != nil {
			return "", false
		}
		raw = parsed.Path
	}

	if raw == "" || strings.Contains(raw, "..") || strings.Contains(raw, "?") || strings.Contains(raw, "#") {
		return "", false
	}

	clean := path.Clean("/" + strings.TrimPrefix(raw, "/"))
	if clean == "/" || clean == "." {
		return "", false
	}

	return clean, true
}

// EditFileName 编辑文件名或者备注
func (b *FileUploadAndDownloadApi) EditFileName(c *gin.Context) {
	if !ensureFileLibraryManager(c) {
		return
	}
	operatorUserID := utils.GetUserID(c)
	operatorAuthorityID := utils.GetUserAuthorityId(c)

	var file example.ExaFileUploadAndDownload
	err := c.ShouldBindJSON(&file)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	err = fileUploadAndDownloadService.EditFileName(file, operatorUserID, operatorAuthorityID)
	if err != nil {
		global.GVA_LOG.Error("编辑失败!", zap.Error(err))
		response.FailWithMessage("编辑失败", c)
		return
	}
	response.OkWithMessage("编辑成功", c)
}

// DeleteFile
// @Tags      ExaFileUploadAndDownload
// @Summary   删除文件
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body      example.ExaFileUploadAndDownload  true  "传入文件里面id即可"
// @Success   200   {object}  response.Response{msg=string}     "删除文件"
// @Router    /fileUploadAndDownload/deleteFile [post]
func (b *FileUploadAndDownloadApi) DeleteFile(c *gin.Context) {
	if !ensureFileLibraryManager(c) {
		return
	}
	operatorUserID := utils.GetUserID(c)
	operatorAuthorityID := utils.GetUserAuthorityId(c)

	var file example.ExaFileUploadAndDownload
	err := c.ShouldBindJSON(&file)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if err := fileUploadAndDownloadService.DeleteFile(file, operatorUserID, operatorAuthorityID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// GetFileList
// @Tags      ExaFileUploadAndDownload
// @Summary   分页文件列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.ExaAttachmentCategorySearch                                        true  "页码, 每页大小, 分类id"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页文件列表,返回包括列表,总数,页码,每页数量"
// @Router    /fileUploadAndDownload/getFileList [post]
func (b *FileUploadAndDownloadApi) GetFileList(c *gin.Context) {
	if !ensureFileLibraryManager(c) {
		return
	}
	operatorUserID := utils.GetUserID(c)
	operatorAuthorityID := utils.GetUserAuthorityId(c)

	var pageInfo request.ExaAttachmentCategorySearch
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	list, total, err := fileUploadAndDownloadService.GetFileRecordInfoList(pageInfo, operatorUserID, operatorAuthorityID)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// ImportURL
// @Tags      ExaFileUploadAndDownload
// @Summary   导入URL
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body      example.ExaFileUploadAndDownload  true  "对象"
// @Success   200   {object}  response.Response{msg=string}     "导入URL"
// @Router    /fileUploadAndDownload/importURL [post]
func (b *FileUploadAndDownloadApi) ImportURL(c *gin.Context) {
	if !ensureFileLibraryManager(c) {
		return
	}
	operatorUserID := utils.GetUserID(c)

	var file []example.ExaFileUploadAndDownload
	err := c.ShouldBindJSON(&file)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if err := fileUploadAndDownloadService.ImportURL(&file, operatorUserID); err != nil {
		global.GVA_LOG.Error("导入URL失败!", zap.Error(err))
		response.FailWithMessage("导入URL失败", c)
		return
	}
	response.OkWithMessage("导入URL成功", c)
}

// SignURL
// @Tags      ExaFileUploadAndDownload
// @Summary   生成防盗链签名URL
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.SignURLRequest  true  "文件路径"
// @Success   200   {object}  response.Response{data=object,msg=string}  "生成成功"
// @Router    /fileUploadAndDownload/signURL [post]
func (b *FileUploadAndDownloadApi) SignURL(c *gin.Context) {
	var req request.SignURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	normalizedPath, ok := normalizeSignFilePath(req.FilePath)
	if !ok {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	if !isFileLibraryManagerRole(utils.GetUserAuthorityId(c)) {
		userID := utils.GetUserID(c)
		if userID == 0 || !isPublicSignMediaPathAllowed(normalizedPath, userID) {
			response.FailWithMessage(i18n.T(c, "noPermission"), c)
			return
		}
	}
	signedURL := upload.SignURL(normalizedPath)
	response.OkWithDetailed(i18n.LocalizeResponseData(c, gin.H{"url": signedURL}), "生成成功", c)
}

// GetHotlinkConfig
// @Tags      ExaFileUploadAndDownload
// @Summary   获取防盗链配置信息
// @Security  ApiKeyAuth
// @Produce   application/json
// @Success   200   {object}  response.Response{data=object,msg=string}  "获取成功"
// @Router    /fileUploadAndDownload/hotlinkConfig [get]
func (b *FileUploadAndDownloadApi) GetHotlinkConfig(c *gin.Context) {
	cfg := global.GVA_CONFIG.Hotlink
	response.OkWithDetailed(gin.H{
		"enabled":   cfg.Enabled,
		"cdnDomain": cfg.CdnDomain,
	}, "获取成功", c)
}

// ListOSSFolders
// @Tags      ExaFileUploadAndDownload
// @Summary   列举OSS存储中的文件夹
// @Security  ApiKeyAuth
// @Produce   application/json
// @Success   200  {object}  response.Response{data=object,msg=string}  "获取成功"
// @Router    /fileUploadAndDownload/listFolders [get]
func (b *FileUploadAndDownloadApi) ListOSSFolders(c *gin.Context) {
	if !ensureFileLibraryManager(c) {
		return
	}

	oss := upload.NewOss()
	folders, err := upload.ListOSSFolders(oss)
	if err != nil {
		global.GVA_LOG.Error("列举文件夹失败!", zap.Error(err))
		response.FailWithMessage("列举文件夹失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"folders": folders}, "获取成功", c)
}
