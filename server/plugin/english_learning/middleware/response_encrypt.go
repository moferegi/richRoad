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
	uniEncryptHeader   = "X-Resp-Encrypt"
	uniEncryptPlatform = "X-Client-Platform"
	uniEncryptSalt     = "uni-api-v1"
)

type encryptWriter struct {
	gin.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (w *encryptWriter) WriteHeader(code int) {
	w.statusCode = code
}

func (w *encryptWriter) Write(data []byte) (int, error) {
	return w.body.Write(data)
}

func LearningResponseEncrypt() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !shouldEncryptForUni(c) {
			c.Next()
			return
		}

		ew := &encryptWriter{
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

		code, ok := envelope["code"].(float64)
		if !ok || int(code) != 0 {
			ew.ResponseWriter.WriteHeader(ew.statusCode)
			_, _ = ew.ResponseWriter.Write(ew.body.Bytes())
			return
		}

		if _, alreadyEnc := envelope["_enc"]; alreadyEnc {
			ew.ResponseWriter.WriteHeader(ew.statusCode)
			_, _ = ew.ResponseWriter.Write(ew.body.Bytes())
			return
		}

		cipherPayload, err := encryptResponseData(c, envelope["data"])
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

func shouldEncryptForUni(c *gin.Context) bool {
	if c == nil || c.Request == nil {
		return false
	}
	if !utils.GetBoolSetting("learning_api_encrypt_enabled", "LEARNING_API_ENCRYPT_ENABLED", true) {
		return false
	}
	if strings.TrimSpace(c.GetHeader(uniEncryptHeader)) != "1" {
		return false
	}
	platform := strings.ToLower(strings.TrimSpace(c.GetHeader(uniEncryptPlatform)))
	return platform == "uni" || platform == "uniapp" || platform == "uni-app"
}

func encryptResponseData(c *gin.Context, data interface{}) (map[string]interface{}, error) {
	plain, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	key := deriveTokenKey(c.GetHeader("x-token"))
	iv := make([]byte, aes.BlockSize)
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	cipherText, err := aesCBCEncrypt(plain, key, iv)
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

func deriveTokenKey(token string) []byte {
	sum := sha256.Sum256([]byte(strings.TrimSpace(token) + "|" + uniEncryptSalt))
	return sum[:]
}

func aesCBCEncrypt(plain []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	padded := pkcs7Pad(plain, aes.BlockSize)
	crypted := make([]byte, len(padded))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(crypted, padded)
	return crypted, nil
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	pad := blockSize - len(data)%blockSize
	if pad == 0 {
		pad = blockSize
	}
	padding := bytes.Repeat([]byte{byte(pad)}, pad)
	return append(data, padding...)
}
