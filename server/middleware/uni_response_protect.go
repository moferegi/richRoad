package middleware

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
)

const (
	uniProtectHeaderEncrypt  = "X-Resp-Encrypt"
	uniProtectHeaderPlatform = "X-Client-Platform"
	uniProtectSalt           = "uni-api-v1"
)

var uniProtectPathPrefixes = []string{
	"/api/banner/getBannerList",
	"/api/kefu/getKefuPublic",
	"/api/cs/config/get",
	"/api/sysConfig/getSysConfigByKey",
	"/api/sysConfig/getLoginConfig",
	"/api/phoneAreaCode/getEnabledPhoneAreaCodes",
	"/api/base/captcha",
	"/api/clientUser/login",
	"/api/clientUser/register",
	"/api/clientUser/phoneLogin",
	"/api/clientUser/phoneRegister",
	"/api/clientUser/getUserInfo",
	"/api/popup/getActivePopups",
}

type uniProtectWriter struct {
	gin.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (w *uniProtectWriter) WriteHeader(code int) {
	w.statusCode = code
}

func (w *uniProtectWriter) Write(data []byte) (int, error) {
	return w.body.Write(data)
}

func UniResponseProtect() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !shouldProtectForUni(c) {
			c.Next()
			return
		}

		ew := &uniProtectWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBuffer(nil),
			statusCode:     http.StatusOK,
		}
		c.Writer = ew
		c.Next()

		if ew.statusCode != http.StatusOK {
			ew.ResponseWriter.WriteHeader(ew.statusCode)
			_, _ = ew.ResponseWriter.Write(ew.body.Bytes())
			return
		}
		if ew.body.Len() == 0 {
			ew.ResponseWriter.WriteHeader(ew.statusCode)
			return
		}

		var envelope map[string]interface{}
		if err := json.Unmarshal(ew.body.Bytes(), &envelope); err != nil {
			ew.ResponseWriter.WriteHeader(ew.statusCode)
			_, _ = ew.ResponseWriter.Write(ew.body.Bytes())
			return
		}

		envelope["data"] = utils.LocalizeI18nPayloadByContext(c, envelope["data"])

		code, ok := envelope["code"].(float64)
		if !ok || int(code) != 0 {
			output, _ := json.Marshal(envelope)
			ew.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
			ew.ResponseWriter.WriteHeader(ew.statusCode)
			_, _ = ew.ResponseWriter.Write(output)
			return
		}

		if !isUniEncryptEnabled() {
			output, _ := json.Marshal(envelope)
			ew.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
			ew.ResponseWriter.WriteHeader(ew.statusCode)
			_, _ = ew.ResponseWriter.Write(output)
			return
		}

		if _, alreadyEnc := envelope["_enc"]; alreadyEnc {
			ew.ResponseWriter.WriteHeader(ew.statusCode)
			_, _ = ew.ResponseWriter.Write(ew.body.Bytes())
			return
		}

		cipherPayload, err := uniProtectEncryptData(c, envelope["data"])
		if err != nil {
			ew.ResponseWriter.WriteHeader(ew.statusCode)
			_, _ = ew.ResponseWriter.Write(ew.body.Bytes())
			return
		}

		envelope["data"] = cipherPayload
		envelope["_enc"] = "aes-cbc"
		envelope["_enc_ver"] = 1

		output, err := json.Marshal(envelope)
		if err != nil {
			ew.ResponseWriter.WriteHeader(ew.statusCode)
			_, _ = ew.ResponseWriter.Write(ew.body.Bytes())
			return
		}

		ew.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
		ew.ResponseWriter.Header().Set("X-Resp-Encrypted", "1")
		ew.ResponseWriter.WriteHeader(ew.statusCode)
		_, _ = ew.ResponseWriter.Write(output)
	}
}

func shouldProtectForUni(c *gin.Context) bool {
	if c == nil || c.Request == nil {
		return false
	}
	if strings.TrimSpace(c.GetHeader(uniProtectHeaderEncrypt)) != "1" {
		return false
	}
	platform := strings.ToLower(strings.TrimSpace(c.GetHeader(uniProtectHeaderPlatform)))
	if platform != "uni" && platform != "uniapp" && platform != "uni-app" {
		return false
	}
	path := strings.ToLower(strings.TrimSpace(c.Request.URL.Path))
	for _, prefix := range uniProtectPathPrefixes {
		if strings.HasPrefix(path, strings.ToLower(prefix)) {
			return true
		}
	}
	return false
}

func isUniEncryptEnabled() bool {
	return utils.GetBoolSetting("learning_api_encrypt_enabled", "LEARNING_API_ENCRYPT_ENABLED", true)
}

func uniProtectEncryptData(c *gin.Context, data interface{}) (map[string]interface{}, error) {
	plain, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	key := uniProtectDeriveTokenKey(c.GetHeader("x-token"))
	iv := make([]byte, aes.BlockSize)
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	cipherText, err := uniProtectAesCBCEncrypt(plain, key, iv)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"alg":     "aes-cbc",
		"payload": base64.StdEncoding.EncodeToString(cipherText),
		"iv":      base64.StdEncoding.EncodeToString(iv),
		"ts":      time.Now().Unix(),
	}, nil
}

func uniProtectDeriveTokenKey(token string) []byte {
	sum := sha256.Sum256([]byte(strings.TrimSpace(token) + "|" + uniProtectSalt))
	return sum[:]
}

func uniProtectAesCBCEncrypt(plain []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	padded := uniProtectPkcs7Pad(plain, aes.BlockSize)
	crypted := make([]byte, len(padded))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(crypted, padded)
	return crypted, nil
}

func uniProtectPkcs7Pad(data []byte, blockSize int) []byte {
	pad := blockSize - len(data)%blockSize
	if pad == 0 {
		pad = blockSize
	}
	padding := bytes.Repeat([]byte{byte(pad)}, pad)
	return append(data, padding...)
}
