package utils

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	jwt "github.com/golang-jwt/jwt/v5"
	goredis "github.com/redis/go-redis/v9"
)

type JWT struct {
	SigningKey []byte
}

var (
	TokenValid            = errors.New("未知错误")
	TokenExpired          = errors.New("token已过期")
	TokenNotValidYet      = errors.New("token尚未激活")
	TokenMalformed        = errors.New("这不是一个token")
	TokenSignatureInvalid = errors.New("无效签名")
	TokenInvalid          = errors.New("无法处理此token")
)

func NewJWT() *JWT {
	return &JWT{
		[]byte(global.GVA_CONFIG.JWT.SigningKey),
	}
}

func (j *JWT) CreateClaims(baseClaims request.BaseClaims) request.CustomClaims {
	bf, _ := ParseDuration(global.GVA_CONFIG.JWT.BufferTime)
	ep, _ := ParseDuration(global.GVA_CONFIG.JWT.ExpiresTime)
	claims := request.CustomClaims{
		BaseClaims: baseClaims,
		BufferTime: int64(bf / time.Second), // 缓冲时间1天 缓冲时间内会获得新的token刷新令牌 此时一个用户会存在两个有效令牌 但是前端只留一个 另一个会丢失
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"GVA"},                   // 受众
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1000)), // 签名生效时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ep)),    // 过期时间 7天  配置文件
			Issuer:    global.GVA_CONFIG.JWT.Issuer,              // 签名的发行者
		},
	}
	return claims
}

// CreateToken 创建一个token
func (j *JWT) CreateToken(claims request.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// CreateTokenByOldToken 旧token 换新token 使用归并回源避免并发问题
func (j *JWT) CreateTokenByOldToken(oldToken string, claims request.CustomClaims) (string, error) {
	v, err, _ := global.GVA_Concurrency_Control.Do("JWT:"+oldToken, func() (interface{}, error) {
		return j.CreateToken(claims)
	})
	return v.(string), err
}

// ParseToken 解析 token
func (j *JWT) ParseToken(tokenString string) (*request.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &request.CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, TokenExpired
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, TokenMalformed
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return nil, TokenSignatureInvalid
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, TokenNotValidYet
		default:
			return nil, TokenInvalid
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*request.CustomClaims); ok && token.Valid {
			return claims, nil
		}
	}
	return nil, TokenValid
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: SetRedisJWT
//@description: jwt存入redis并设置过期时间
//@param: jwt string, userName string
//@return: err error

func SetRedisJWT(jwt string, userName string) (err error) {
	// 此处过期时间等于jwt过期时间
	dr, err := ParseDuration(global.GVA_CONFIG.JWT.ExpiresTime)
	if err != nil {
		return err
	}
	timer := dr
	err = global.GVA_REDIS.Set(context.Background(), userName, jwt, timer).Err()
	return err
}

// DeviceTokenKey 返回用户设备令牌 Redis SortedSet key
func DeviceTokenKey(username string) string {
	return "jwt_devices:" + username
}

// AddDeviceToken 将新 JWT 加入用户设备列表，并强制执行设备数限制。
// 返回被踢出（最早登录）的令牌列表，调用方应将它们加入黑名单。
func AddDeviceToken(username, token string, expiry time.Duration, maxDevices int) (evicted []string, err error) {
	if global.GVA_REDIS == nil {
		return nil, nil
	}
	ctx := context.Background()
	key := DeviceTokenKey(username)
	now := time.Now()
	score := float64(now.Add(expiry).Unix())

	// 先清除已过期的令牌
	_ = global.GVA_REDIS.ZRemRangeByScore(ctx, key, "-inf", strconv.FormatInt(now.Unix(), 10)).Err()

	if maxDevices > 0 {
		count, _ := global.GVA_REDIS.ZCard(ctx, key).Result()
		excess := int(count) - maxDevices + 1 // +1 为新令牌计算需要踢出的数量
		if excess > 0 {
			oldest, _ := global.GVA_REDIS.ZRange(ctx, key, 0, int64(excess-1)).Result()
			evicted = oldest
			_ = global.GVA_REDIS.ZRemRangeByRank(ctx, key, 0, int64(excess-1)).Err()
		}
	}

	err = global.GVA_REDIS.ZAdd(ctx, key, goredis.Z{
		Score:  score,
		Member: token,
	}).Err()
	if err != nil {
		return evicted, err
	}
	_ = global.GVA_REDIS.Expire(ctx, key, expiry).Err()
	return evicted, nil
}

// RemoveDeviceToken 从用户设备列表移除指定令牌
func RemoveDeviceToken(username, token string) {
	if global.GVA_REDIS == nil {
		return
	}
	_ = global.GVA_REDIS.ZRem(context.Background(), DeviceTokenKey(username), token).Err()
}

// RenewDeviceToken 用新令牌替换旧令牌，不改变设备数量（续签时使用）
func RenewDeviceToken(username, oldToken, newToken string, expiry time.Duration) {
	if global.GVA_REDIS == nil {
		return
	}
	ctx := context.Background()
	key := DeviceTokenKey(username)
	_ = global.GVA_REDIS.ZRem(ctx, key, oldToken).Err()
	_ = global.GVA_REDIS.ZAdd(ctx, key, goredis.Z{
		Score:  float64(time.Now().Add(expiry).Unix()),
		Member: newToken,
	}).Err()
	_ = global.GVA_REDIS.Expire(ctx, key, expiry).Err()
}

// GetDeviceCount 返回用户当前活跃的登录设备数量
func GetDeviceCount(username string) int64 {
	if global.GVA_REDIS == nil {
		return 0
	}
	ctx := context.Background()
	key := DeviceTokenKey(username)
	_ = global.GVA_REDIS.ZRemRangeByScore(ctx, key, "-inf", strconv.FormatInt(time.Now().Unix(), 10)).Err()
	count, _ := global.GVA_REDIS.ZCard(ctx, key).Result()
	return count
}
