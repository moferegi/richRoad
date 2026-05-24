package api

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	wxGlobal "github.com/flipped-aurora/gin-vue-admin/server/plugin/wxpay/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/wxpay/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/wxpay/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"go.uber.org/zap"
)

type WxpayApi struct{}

const (
	wechatpayTimestampToleranceSeconds = 300
	wechatpayNonceTTL                  = 10 * time.Minute
)

type notifyNonceCache struct {
	mu    sync.Mutex
	items map[string]time.Time
}

var wxpayNotifyNonceCache = &notifyNonceCache{items: map[string]time.Time{}}

func (c *notifyNonceCache) seenOrStore(nonce string, ttl time.Duration) bool {
	now := time.Now()
	expiresAt := now.Add(ttl)

	c.mu.Lock()
	defer c.mu.Unlock()

	for key, exp := range c.items {
		if exp.Before(now) {
			delete(c.items, key)
		}
	}

	if exp, ok := c.items[nonce]; ok && exp.After(now) {
		return true
	}
	c.items[nonce] = expiresAt
	return false
}

func buildNotifyHandler(ctx context.Context) (*notify.Handler, error) {
	conf := wxGlobal.GlobalConfig
	if strings.TrimSpace(conf.MchID) == "" || strings.TrimSpace(conf.MchAPIv3Key) == "" || strings.TrimSpace(conf.MchCertificateSerialNumber) == "" || strings.TrimSpace(conf.KeyPath) == "" {
		return nil, errors.New("wxpay callback verification config incomplete")
	}

	privateKey, err := loadRSAPrivateKey(conf.KeyPath)
	if err != nil {
		return nil, err
	}

	certificateDownloader, err := downloader.NewCertificateDownloader(ctx, conf.MchID, privateKey, conf.MchCertificateSerialNumber, conf.MchAPIv3Key)
	if err != nil {
		return nil, err
	}

	verifier := verifiers.NewSHA256WithRSAVerifier(certificateDownloader)
	return notify.NewRSANotifyHandler(conf.MchAPIv3Key, verifier)
}

func loadRSAPrivateKey(privateKeyPath string) (*rsa.PrivateKey, error) {
	keyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, errors.New("invalid private key pem")
	}

	if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		rsaKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("private key is not rsa")
		}
		return rsaKey, nil
	}

	rsaKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsaKey, nil
}

func validateNotifyBasicReplayGuard(c *gin.Context) (string, error) {
	timestamp := strings.TrimSpace(c.GetHeader("Wechatpay-Timestamp"))
	nonce := strings.TrimSpace(c.GetHeader("Wechatpay-Nonce"))
	signature := strings.TrimSpace(c.GetHeader("Wechatpay-Signature"))
	serial := strings.TrimSpace(c.GetHeader("Wechatpay-Serial"))

	if timestamp == "" || nonce == "" || signature == "" || serial == "" {
		return "", errors.New("missing wechatpay signature headers")
	}

	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return "", errors.New("invalid wechatpay timestamp")
	}

	delta := math.Abs(float64(time.Now().Unix() - ts))
	if delta > wechatpayTimestampToleranceSeconds {
		return "", errors.New("wechatpay callback timestamp expired")
	}

	return nonce, nil
}

func markNotifyNonceUsed(nonce string) error {
	if wxpayNotifyNonceCache.seenOrStore(nonce, wechatpayNonceTTL) {
		return errors.New("wechatpay callback nonce replayed")
	}
	return nil
}

// @Tags Wxpay
// @Summary 获取微信支付二维码和ID
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"发送成功"}"
// @Router /wxpay/getPayCode[post]
func (p *WxpayApi) GetPayCode(c *gin.Context) {
	var order model.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	order.CustomerID = utils.GetUserID(c)
	if order.CustomerID == 0 {
		response.FailWithMessage("请先登录", c)
		return
	}
	if err, codeUrl, codeId := service.ServiceGroupApp.GetPayCode(order); err != nil {
		global.GVA_LOG.Error("失败!", zap.Error(err))
		response.FailWithMessage("获取付款码失败", c)
	} else {
		response.OkWithData(gin.H{
			"codeUrl": codeUrl,
			"codeId":  codeId,
		}, c)
	}
}

func (p *WxpayApi) CheckNeedPay(c *gin.Context) {
	var order model.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	order.CustomerID = utils.GetUserID(c)
	if order.CustomerID == 0 {
		response.FailWithMessage("请先登录", c)
		return
	}

	err, needPay := service.ServiceGroupApp.CheckNeedPay(order)

	if err != nil {
		global.GVA_LOG.Error("失败!", zap.Error(err))
		response.FailWithMessage("检查订单失败", c)
		return
	}
	response.OkWithData(needPay, c)
}

// @Tags Wxpay
// @Summary 获取微信支付二维码和ID
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"发送成功"}"
// @Router /wxpay/GetPayParams[post]
func (p *WxpayApi) GetPayParams(c *gin.Context) {
	var order model.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	order.CustomerID = utils.GetUserID(c)
	if order.CustomerID == 0 {
		response.FailWithMessage("请先登录", c)
		return
	}
	if err, patConf := service.ServiceGroupApp.GetPayConf(order); err != nil {
		global.GVA_LOG.Error("失败!", zap.Error(err))
		response.FailWithMessage("获取付款码失败", c)
	} else {
		response.OkWithData(patConf, c)
	}
}

// @Tags Wxpay
// @Summary 获取支付结果
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"发送成功"}"
// @Router /wxpay/getOrderById[get]
func (p *WxpayApi) GetOrderById(c *gin.Context) {
	id := c.Query("orderID")
	userID := utils.GetUserID(c)
	if userID == 0 {
		response.FailWithMessage("请先登录", c)
		return
	}
	if err, data := service.ServiceGroupApp.GetOrderById(id, userID); err != nil {
		global.GVA_LOG.Error("失败!", zap.Error(err))
		response.FailWithMessage("获取付款码失败", c)
	} else {
		response.OkWithData(data, c)
	}
}

// @Tags Wxpay
// @Summary 回调支付结果
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"发送成功"}"
// @Router /wxpay/payAction[post]
func (p *WxpayApi) PayAction(c *gin.Context) {
	nonce, err := validateNotifyBasicReplayGuard(c)
	if err != nil {
		global.GVA_LOG.Error("wxpay callback headers validation failed", zap.Error(err))
		c.JSON(200, gin.H{
			"code":    "FAIL",
			"message": "失败",
		})
		return
	}

	handler, err := buildNotifyHandler(c.Request.Context())
	if err != nil {
		global.GVA_LOG.Error("wxpay callback verifier init failed", zap.Error(err))
		c.JSON(200, gin.H{
			"code":    "FAIL",
			"message": "失败",
		})
		return
	}

	var payOrder model.PayOrder
	if _, err := handler.ParseNotifyRequest(c.Request.Context(), c.Request, &payOrder); err != nil {
		global.GVA_LOG.Error("wxpay callback signature parse failed", zap.Error(err))
		c.JSON(200, gin.H{
			"code":    "FAIL",
			"message": "失败",
		})
		return
	}

	if err := markNotifyNonceUsed(nonce); err != nil {
		global.GVA_LOG.Error("wxpay callback replay validation failed", zap.Error(err))
		c.JSON(200, gin.H{
			"code":    "FAIL",
			"message": "失败",
		})
		return
	}

	if err := service.ServiceGroupApp.PayActionOrder(payOrder); err != nil {
		global.GVA_LOG.Error("失败!", zap.Error(err))
		c.JSON(200, gin.H{
			"code":    "FAIL",
			"message": "失败",
		})
	} else {
		c.JSON(200, gin.H{
			"code":    "SUCCESS",
			"message": "接收成功",
		})
	}
}
