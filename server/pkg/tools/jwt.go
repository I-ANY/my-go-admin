package tools

import (
	"biz-auto-api/pkg/consts"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"time"
)

// CreateToken 生成Token：
func CreateToken(id, issuer string, periodMinutes int) (tokenString string, err error) {
	secretKey := []byte(consts.SecretKey)
	m := time.Duration(periodMinutes)
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * m)),
		Issuer:    issuer,
		ID:        id,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err = token.SignedString(secretKey)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return
}

// ParseToken 解析Token
func ParseToken(tokenSrt string) (claims *jwt.RegisteredClaims, err error) {
	secretKey := []byte(consts.SecretKey)
	var token *jwt.Token
	token, err = jwt.ParseWithClaims(tokenSrt, &jwt.RegisteredClaims{}, func(*jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	claims = token.Claims.(*jwt.RegisteredClaims)
	return
}
