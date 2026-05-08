package api

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"sort"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/model"
	csReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/service"
	gvaService "github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MessageApi struct{}

const (
	defaultUploadMaxSizeMB = 5
	maxUploadMaxSizeMB     = 50
)

var imageMimeByExt = map[string]string{
	"jpg":  "image/jpeg",
	"jpeg": "image/jpeg",
	"png":  "image/png",
	"webp": "image/webp",
	"gif":  "image/gif",
}

var (
	errImageTooLarge = errors.New("image too large")
	errImageType     = errors.New("image type invalid")
	errImageEmpty    = errors.New("image empty")
)

func parseAllowImageExt(raw string) ([]string, map[string]struct{}) {
	parts := strings.Split(raw, ",")
	list := make([]string, 0, len(parts))
	set := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		ext := strings.ToLower(strings.TrimSpace(part))
		ext = strings.TrimPrefix(ext, ".")
		if ext == "" {
			continue
		}
		if _, ok := imageMimeByExt[ext]; !ok {
			continue
		}
		if _, ok := set[ext]; ok {
			continue
		}
		set[ext] = struct{}{}
		list = append(list, ext)
	}
	if len(list) == 0 {
		return parseAllowImageExt("jpg,jpeg,png,webp,gif")
	}
	sort.Strings(list)
	return list, set
}

func buildAllowImageMime(extSet map[string]struct{}) map[string]struct{} {
	mimes := make(map[string]struct{}, len(extSet))
	for ext := range extSet {
		if mime, ok := imageMimeByExt[ext]; ok {
			mimes[mime] = struct{}{}
		}
	}
	if len(mimes) == 0 {
		for _, mime := range imageMimeByExt {
			mimes[mime] = struct{}{}
		}
	}
	return mimes
}

func normalizeUploadMaxSizeMB(raw int) int {
	if raw <= 0 {
		return defaultUploadMaxSizeMB
	}
	if raw > maxUploadMaxSizeMB {
		return maxUploadMaxSizeMB
	}
	return raw
}

func validateChatImageHeader(header *multipart.FileHeader, maxSize int64, allowImageExt map[string]struct{}, allowImageMime map[string]struct{}) error {
	if header == nil {
		return errImageEmpty
	}
	if header.Size <= 0 {
		return errImageEmpty
	}
	if header.Size > maxSize {
		return errImageTooLarge
	}

	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(header.Filename)), ".")
	if _, ok := allowImageExt[ext]; !ok {
		return errImageType
	}

	f, err := header.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	buf := make([]byte, 512)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return err
	}
	if n == 0 {
		return errImageEmpty
	}

	if _, ok := allowImageMime[http.DetectContentType(buf[:n])]; !ok {
		return errImageType
	}
	return nil
}

// UploadImage 上传聊天图片（服务端强校验）
// @Tags CustomerService
// @Summary 上传聊天图片
// @Security ApiKeyAuth
// @accept multipart/form-data
// @Produce application/json
// @Param file formData file true "图片文件(格式和大小由客服配置决定)"
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "上传成功"
// @Router /cs/message/upload [post]
func (a *MessageApi) UploadImage(c *gin.Context) {
	cfg, err := service.Service.ConfigService.GetConfig()
	if err != nil {
		response.FailWithMessage("读取上传配置失败", c)
		return
	}

	maxSizeMB := normalizeUploadMaxSizeMB(cfg.UploadMaxSizeMB)
	maxUploadImageSize := int64(maxSizeMB) * 1024 * 1024
	allowExtList, allowExtSet := parseAllowImageExt(cfg.UploadAllowExt)
	allowMimeSet := buildAllowImageMime(allowExtSet)

	// 提前限制请求体，防止超大表单占用内存
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadImageSize+256*1024)
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "too large") {
			response.FailWithMessage(fmt.Sprintf("图片大小不能超过 %dMB", maxSizeMB), c)
			return
		}
		response.FailWithMessage("接收图片失败", c)
		return
	}

	authID := utils.GetUserAuthorityId(c)
	if authID != 8080 && authID != 888 {
		if _, err := service.Service.AgentService.GetEnabledByUserID(utils.GetUserID(c)); err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
	}

	if err := validateChatImageHeader(header, maxUploadImageSize, allowExtSet, allowMimeSet); err != nil {
		switch err {
		case errImageTooLarge:
			response.FailWithMessage(fmt.Sprintf("图片大小不能超过 %dMB", maxSizeMB), c)
		case errImageType:
			response.FailWithMessage(fmt.Sprintf("仅支持 %s 图片", strings.ToUpper(strings.Join(allowExtList, "/"))), c)
		default:
			response.FailWithMessage("图片格式不合法", c)
		}
		return
	}

	file, err := gvaService.ServiceGroupApp.ExampleServiceGroup.FileUploadAndDownloadService.UploadFile(header, "0", 0, "cloth-on/kefu", "", "kefu")
	if err != nil {
		global.GVA_LOG.Error("客服图片上传失败", zap.Error(err))
		response.FailWithMessage("上传失败", c)
		return
	}

	response.OkWithDetailed(gin.H{
		"file": gin.H{"url": file.Url},
		"url":  file.Url,
	}, "上传成功", c)
}

// GetMessageHistory 获取消息历史
// @Tags CustomerService
// @Summary 获取会话消息历史（分页）
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query csReq.MessageSearch true "分页参数"
// @Success 200 {object} response.Response{data=response.PageResult} "获取成功"
// @Router /cs/message/history [get]
func (a *MessageApi) GetMessageHistory(c *gin.Context) {
	var search csReq.MessageSearch
	if err := c.ShouldBindQuery(&search); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userID := utils.GetUserID(c)
	authID := utils.GetUserAuthorityId(c)

	if authID == 8080 {
		if err := service.Service.ConversationService.CheckClientOwnership(userID, search.ConversationID); err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
	} else if authID != 888 {
		if _, err := service.Service.AgentService.GetEnabledByUserID(userID); err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
		if err := service.Service.ConversationService.CheckAgentAccess(userID, search.ConversationID, true); err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
	}

	list, total, err := service.Service.MessageService.GetHistory(search)
	if err != nil {
		global.GVA_LOG.Error("获取消息历史失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     search.Page,
		PageSize: search.PageSize,
	}, "获取成功", c)
}

// SendMessage 通过 REST 发送消息（坐席后台补发等场景）
// @Tags CustomerService
// @Summary 发送消息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body csReq.SendMessageReq true "消息内容"
// @Success 200 {object} response.Response{msg=string} "发送成功"
// @Router /cs/agent/message/send [post]
func (a *MessageApi) SendMessage(c *gin.Context) {
	var req csReq.SendMessageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	senderID := utils.GetUserID(c)
	authID := utils.GetUserAuthorityId(c)
	if authID == 8080 {
		response.FailWithMessage("客户端用户不能调用坐席发送接口", c)
		return
	}
	if authID != 888 {
		if _, err := service.Service.AgentService.GetEnabledByUserID(senderID); err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
		if err := service.Service.ConversationService.CheckAgentAccess(senderID, req.ConversationID, true); err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
	}
	// 默认 agent 类型（REST 接口一般坐席调用）
	msg, err := service.Service.MessageService.Send(
		model.SenderTypeAgent, senderID, req.ConversationID,
		req.MsgType, req.Content, req.ClientMsgID,
	)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(msg, "发送成功", c)
}

// RevokeMessage 撤回消息
// @Tags CustomerService
// @Summary 撤回消息（2分钟内）
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body csReq.RevokeMessageReq true "消息ID"
// @Success 200 {object} response.Response{msg=string} "撤回成功"
// @Router /cs/message/revoke [post]
func (a *MessageApi) RevokeMessage(c *gin.Context) {
	var req csReq.RevokeMessageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	senderID := utils.GetUserID(c)
	// 从 authority 区分发送方类型：8080 = 客户端用户，其余 = 坐席
	authID := utils.GetUserAuthorityId(c)
	senderType := "agent"
	if authID == 8080 {
		senderType = "user"
	} else if authID != 888 {
		if _, err := service.Service.AgentService.GetEnabledByUserID(senderID); err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
	}
	if err := service.Service.MessageService.Revoke(req.MessageID, senderID, senderType); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("撤回成功", c)
}
