package utils

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/payment"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/wxpay/global"
	"os"
	"strings"
)

func CreateClientAndCtx() (error, context.Context, *payment.Payment) {
	PaymentService, err := payment.NewPayment(&payment.UserConfig{
		AppID:       global.GlobalConfig.AppID,                      // 小程序、公众号或者企业微信的appid
		MchID:       global.GlobalConfig.MchID,                      // 商户号 appID
		MchApiV3Key: global.GlobalConfig.MchAPIv3Key,                // 微信V3接口调用必填
		Key:         global.GlobalConfig.MchAPIv2Key,                // 微信V2接口调用必填
		CertPath:    global.GlobalConfig.CertPath,                   // 商户后台支付的Cert证书路径
		KeyPath:     global.GlobalConfig.KeyPath,                    // 商户后台支付的Key证书路径
		SerialNo:    global.GlobalConfig.MchCertificateSerialNumber, // 商户支付证书序列号
		NotifyURL:   global.GlobalConfig.NotifyUrl,
		HttpDebug:   false, // 订单模式不支持debug 请勿打开
		Log: payment.Log{
			Level: "debug",
			// 可以重定向到你的目录下，如果设置File和Error，默认会在当前目录下的wechat文件夹下生成日志
			File:   "info.log",
			Error:  "error.log",
			Stdout: false, //  是否打印在终端
		},
		Http: payment.Http{
			Timeout: 30.0,
			BaseURI: "https://api.mch.weixin.qq.com",
		},
	})
	return err, context.Background(), PaymentService
}

func GeneratePaySign(privateKeyPem string, appId, timeStamp, nonceStr, packageValue string) (string, error) {
	// 解码私钥
	pbyte, err := os.ReadFile(privateKeyPem)
	if err != nil {
		return "", fmt.Errorf("failed to read private key: %v", err)
	}
	block, _ := pem.Decode(pbyte)
	if block == nil {
		return "", fmt.Errorf("failed to decode private key")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %v", err)
	}
	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return "", fmt.Errorf("not an RSA private key")
	}
	// 构造待签名的数据
	data := strings.Join([]string{appId, timeStamp, nonceStr, packageValue}, "\n") + "\n"

	// 计算SHA256哈希
	hash := sha256.Sum256([]byte(data))

	// 使用私钥进行RSA签名
	signature, err := rsa.SignPKCS1v15(rand.Reader, rsaPrivateKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", fmt.Errorf("failed to sign data: %v", err)
	}

	// 将签名转换为Base64编码的字符串
	signatureStr := base64.StdEncoding.EncodeToString(signature)

	return signatureStr, nil
}
