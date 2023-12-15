package jwt

import (
	"errors"
	"hall-server/internal/errx"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
		return "", 0, errx.GenTokenError.E(err)
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
			return nil, errx.TokenMalformed.N()
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, errx.TokenExpired.N()
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, errx.TokenNotValidYet.N()
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return nil, errx.TokenSignatureInvalid.N()
		default:
			return nil, errx.TokenInvalid.E(err)
		}
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errx.TokenInvalid.N()
}
