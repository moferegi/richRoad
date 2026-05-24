package system

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/i18n"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SysErrorApi struct{}

func isSysErrorManageRole(authorityID uint) bool {
	return authorityID == 888 || authorityID == 8881
}

const (
	sysErrorCreateRateLimitConfigKey = "security_sys_error_create_rate_limit_per_minute"
	sysErrorCreateRateLimitEnvKey    = "CS_SYS_ERROR_CREATE_RATE_LIMIT_PER_MINUTE"

	sysErrorCreateRateWindowConfigKey = "security_sys_error_create_rate_limit_window_seconds"
	sysErrorCreateRateWindowEnvKey    = "CS_SYS_ERROR_CREATE_RATE_LIMIT_WINDOW_SECONDS"

	defaultSysErrorCreateRateLimitPerMinute int64 = 30
	defaultSysErrorCreateRateWindowSeconds  int64 = 60
)

type localSysErrorRateCounter struct {
	Count     int64
	ExpiresAt time.Time
}

var localSysErrorCreateRateLimiter = struct {
	mu       sync.Mutex
	counters map[string]localSysErrorRateCounter
}{
	counters: make(map[string]localSysErrorRateCounter),
}

// CreateSysError 创建错误日志
// @Tags SysError
// @Summary 创建错误日志
// @Accept application/json
// @Produce application/json
// @Param data body systemReq.CreateSysErrorRequest true "创建错误日志"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /sysError/createSysError [post]
func (sysErrorApi *SysErrorApi) CreateSysError(c *gin.Context) {
	if allowed, waitSeconds := checkSysErrorCreateRateLimit(c); !allowed {
		bannedIPService.RecordAttack(c.ClientIP(), "sys_error_rate_limit", "", "错误上报频率超限")
		response.FailWithMessage(i18n.T(c, "requestTooFrequent")+" "+strconv.Itoa(waitSeconds)+"s", c)
		return
	}

	// 创建业务用Context
	ctx := c.Request.Context()

	var req systemReq.CreateSysErrorRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}

	form := strings.TrimSpace(req.Form)
	info := strings.TrimSpace(req.Info)
	if form == "" || info == "" {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}

	level := normalizeSysErrorLevel(req.Level)
	sysError := system.SysError{
		Form:   &form,
		Info:   &info,
		Level:  level,
		Status: "未处理",
	}

	err = sysErrorService.CreateSysError(ctx, &sysError)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "createFail"), c)
		return
	}
	response.OkWithMessage(i18n.T(c, "createSuccess"), c)
}

// DeleteSysError 删除错误日志
// @Tags SysError
// @Summary 删除错误日志
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body system.SysError true "删除错误日志"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /sysError/deleteSysError [delete]
func (sysErrorApi *SysErrorApi) DeleteSysError(c *gin.Context) {
	if !isSysErrorManageRole(utils.GetUserAuthorityId(c)) {
		response.FailWithMessage(i18n.T(c, "noPermission"), c)
		return
	}
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	err := sysErrorService.DeleteSysError(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "deleteFail"), c)
		return
	}
	response.OkWithMessage(i18n.T(c, "deleteSuccess"), c)
}

// DeleteSysErrorByIds 批量删除错误日志
// @Tags SysError
// @Summary 批量删除错误日志
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /sysError/deleteSysErrorByIds [delete]
func (sysErrorApi *SysErrorApi) DeleteSysErrorByIds(c *gin.Context) {
	if !isSysErrorManageRole(utils.GetUserAuthorityId(c)) {
		response.FailWithMessage(i18n.T(c, "noPermission"), c)
		return
	}
	// 创建业务用Context
	ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := sysErrorService.DeleteSysErrorByIds(ctx, IDs)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "batchDeleteFail"), c)
		return
	}
	response.OkWithMessage(i18n.T(c, "batchDeleteSuccess"), c)
}

// UpdateSysError 更新错误日志
// @Tags SysError
// @Summary 更新错误日志
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body system.SysError true "更新错误日志"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /sysError/updateSysError [put]
func (sysErrorApi *SysErrorApi) UpdateSysError(c *gin.Context) {
	if !isSysErrorManageRole(utils.GetUserAuthorityId(c)) {
		response.FailWithMessage(i18n.T(c, "noPermission"), c)
		return
	}
	// 从ctx获取标准context进行业务行为
	ctx := c.Request.Context()

	var sysError system.SysError
	err := c.ShouldBindJSON(&sysError)
	if err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	err = sysErrorService.UpdateSysError(ctx, sysError)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "updateFail"), c)
		return
	}
	response.OkWithMessage(i18n.T(c, "updateSuccess"), c)
}

// FindSysError 用id查询错误日志
// @Tags SysError
// @Summary 用id查询错误日志
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询错误日志"
// @Success 200 {object} response.Response{data=system.SysError,msg=string} "查询成功"
// @Router /sysError/findSysError [get]
func (sysErrorApi *SysErrorApi) FindSysError(c *gin.Context) {
	if !isSysErrorManageRole(utils.GetUserAuthorityId(c)) {
		response.FailWithMessage(i18n.T(c, "noPermission"), c)
		return
	}
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	resysError, err := sysErrorService.GetSysError(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "queryFail"), c)
		return
	}
	response.OkWithData(resysError, c)
}

// GetSysErrorList 分页获取错误日志列表
// @Tags SysError
// @Summary 分页获取错误日志列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query systemReq.SysErrorSearch true "分页获取错误日志列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /sysError/getSysErrorList [get]
func (sysErrorApi *SysErrorApi) GetSysErrorList(c *gin.Context) {
	if !isSysErrorManageRole(utils.GetUserAuthorityId(c)) {
		response.FailWithMessage(i18n.T(c, "noPermission"), c)
		return
	}
	// 创建业务用Context
	ctx := c.Request.Context()

	var pageInfo systemReq.SysErrorSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	list, total, err := sysErrorService.GetSysErrorInfoList(ctx, pageInfo)
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

// GetSysErrorSolution 触发错误日志的异步处理
// @Tags SysError
// @Summary 根据ID触发处理：标记为处理中，1分钟后自动改为处理完成
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id query string true "错误日志ID"
// @Success 200 {object} response.Response{msg=string} "处理已提交"
// @Router /sysError/getSysErrorSolution [get]
func (sysErrorApi *SysErrorApi) GetSysErrorSolution(c *gin.Context) {
	if !isSysErrorManageRole(utils.GetUserAuthorityId(c)) {
		response.FailWithMessage(i18n.T(c, "noPermission"), c)
		return
	}
	// 创建业务用Context
	ctx := c.Request.Context()

	// 兼容 id 与 ID 两种参数
	ID := c.Query("id")
	if ID == "" {
		response.FailWithMessage(i18n.T(c, "invalidID"), c)
		return
	}

	err := sysErrorService.GetSysErrorSolution(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("处理触发失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "fail"), c)
		return
	}

	response.OkWithMessage(i18n.T(c, "submitSuccess"), c)
}

func normalizeSysErrorLevel(raw string) string {
	level := strings.ToLower(strings.TrimSpace(raw))
	switch level {
	case "error", "warn", "info", "debug":
		return level
	case "warning":
		return "warn"
	default:
		return "error"
	}
}

func checkSysErrorCreateRateLimit(c *gin.Context) (allowed bool, waitSeconds int) {
	ip := strings.TrimSpace(c.ClientIP())
	if ip == "" {
		ip = "unknown"
	}

	limit := utils.GetInt64Setting(sysErrorCreateRateLimitConfigKey, sysErrorCreateRateLimitEnvKey, defaultSysErrorCreateRateLimitPerMinute)
	if limit <= 0 {
		return true, 0
	}

	windowSeconds := utils.GetInt64Setting(sysErrorCreateRateWindowConfigKey, sysErrorCreateRateWindowEnvKey, defaultSysErrorCreateRateWindowSeconds)
	if windowSeconds <= 0 {
		windowSeconds = defaultSysErrorCreateRateWindowSeconds
	}
	window := time.Duration(windowSeconds) * time.Second

	if global.GVA_REDIS == nil {
		return checkSysErrorCreateRateLimitLocal(ip, limit, window)
	}

	key := fmt.Sprintf("rate:sysError:create:ip:%s", ip)
	count, err := global.GVA_REDIS.Incr(c.Request.Context(), key).Result()
	if err != nil {
		return checkSysErrorCreateRateLimitLocal(ip, limit, window)
	}
	if count == 1 {
		_ = global.GVA_REDIS.Expire(c.Request.Context(), key, window).Err()
	}
	if count > limit {
		ttl, ttlErr := global.GVA_REDIS.TTL(c.Request.Context(), key).Result()
		if ttlErr != nil || ttl <= 0 {
			return false, int(window.Seconds())
		}
		wait := int(ttl.Seconds())
		if wait < 1 {
			wait = 1
		}
		return false, wait
	}

	return true, 0
}

func checkSysErrorCreateRateLimitLocal(ip string, limit int64, window time.Duration) (allowed bool, waitSeconds int) {
	if limit <= 0 {
		return true, 0
	}
	if window <= 0 {
		window = time.Duration(defaultSysErrorCreateRateWindowSeconds) * time.Second
	}

	now := time.Now()
	localSysErrorCreateRateLimiter.mu.Lock()
	defer localSysErrorCreateRateLimiter.mu.Unlock()

	if len(localSysErrorCreateRateLimiter.counters) > 10000 {
		for key, counter := range localSysErrorCreateRateLimiter.counters {
			if now.After(counter.ExpiresAt) {
				delete(localSysErrorCreateRateLimiter.counters, key)
			}
		}
	}

	counter, ok := localSysErrorCreateRateLimiter.counters[ip]
	if !ok || now.After(counter.ExpiresAt) {
		localSysErrorCreateRateLimiter.counters[ip] = localSysErrorRateCounter{
			Count:     1,
			ExpiresAt: now.Add(window),
		}
		return true, 0
	}

	if counter.Count >= limit {
		remain := int(time.Until(counter.ExpiresAt).Seconds())
		if remain < 1 {
			remain = 1
		}
		return false, remain
	}

	counter.Count++
	localSysErrorCreateRateLimiter.counters[ip] = counter
	return true, 0
}
