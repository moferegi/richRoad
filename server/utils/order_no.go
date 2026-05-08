package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

const (
	orderNoDigits  = "0123456789"
	orderNoLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// GenerateBusinessOrderNo 生成业务订单号，格式：前缀(2位)+5位数字+2位大写字母。
// 示例：sp12345AB / sy54321ZX
func GenerateBusinessOrderNo(prefix string) (string, error) {
	prefix = strings.ToLower(strings.TrimSpace(prefix))
	if len(prefix) != 2 {
		return "", errors.New("订单号前缀长度必须为2位")
	}

	digits, err := randomChars(orderNoDigits, 5)
	if err != nil {
		return "", err
	}
	letters, err := randomChars(orderNoLetters, 2)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%s%s", prefix, digits, letters), nil
}

func randomChars(charset string, length int) (string, error) {
	if length <= 0 {
		return "", errors.New("随机长度必须大于0")
	}
	if len(charset) == 0 {
		return "", errors.New("字符集不能为空")
	}

	b := make([]byte, length)
	max := big.NewInt(int64(len(charset)))
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		b[i] = charset[n.Int64()]
	}
	return string(b), nil
}
