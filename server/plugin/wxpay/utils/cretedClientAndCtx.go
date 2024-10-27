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
	global2 "github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/wxpay/global"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"log"
	"os"
	"strings"
)

func CreateClientAndCtx() (error, context.Context, *core.Client) {
	var (
		mchID                      string = global.GlobalConfig.MchID                      // 商户号
		mchCertificateSerialNumber string = global.GlobalConfig.MchCertificateSerialNumber // 商户证书序列号
		mchAPIv3Key                string = global.GlobalConfig.MchAPIv3Key                // 商户APIv3密钥
	)
	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(global.GlobalConfig.PemPath)
	if err != nil {
		log.Fatal("load merchant private key error")
	}
	ctx := context.Background()
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(mchID, mchCertificateSerialNumber, mchPrivateKey, mchAPIv3Key),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		e := err.Error()
		global2.GVA_LOG.Error("new wechat pay client err:" + e)
		return err, ctx, client
	}
	return err, ctx, client
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
