package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	GenTokenError         = errors.New("生成令牌失败")
	TokenMalformed        = errors.New("令牌格式错误")
	TokenExpired          = errors.New("令牌已过期")
	TokenNotValidYet      = errors.New("令牌尚未激活")
	TokenSignatureInvalid = errors.New("令牌签名无效")
	TokenInvalid          = errors.New("令牌无效")
)

type Invoker interface {
	Generate(bundleName, deviceType, pid, openId string) (string, int64, error)
	Validate(tokenString string) (*Claims, error)
}

func (c *Component) Generate(bundleName, deviceType, pid, openId string) (string, int64, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(c.config.ExpireDuration) * time.Hour)
	claims := &Claims{
		Pid:        pid,
		OpenId:     openId,
		DeviceType: deviceType,
		BundleName: bundleName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	sign, err := token.SignedString(c.config.SecretKey)
	if err != nil {
		return "", 0, GenTokenError
	}
	return sign, expiresAt.Unix(), nil
}

func (c *Component) Validate(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return c.config.SecretKey, nil
	})
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, TokenMalformed
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, TokenExpired
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, TokenNotValidYet
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return nil, TokenSignatureInvalid
		default:
			return nil, TokenInvalid
		}
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}
