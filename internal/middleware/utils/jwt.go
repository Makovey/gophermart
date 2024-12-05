package utils

import (
	"errors"
	"fmt"
	"github.com/Makovey/gophermart/internal/logger"
	"github.com/golang-jwt/jwt/v5"
	"os"
)

var (
	ErrParseToken    = errors.New("failed to parse token")
	ErrSigningMethod = errors.New("unexpected signing method")
	ErrInvalidToken  = errors.New("invalid token")
	ErrTokenExpired  = errors.New("token is expired")
)

type JWTUtils struct {
	log logger.Logger
}

func NewJWTUtils(logger logger.Logger) JWTUtils {
	return JWTUtils{log: logger}
}

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

func (j JWTUtils) ParseUserID(tokenString string) (string, error) {
	fn := "utils.jwt.parseToken"

	var claims Claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			j.log.Warn(fmt.Sprintf("%s: unexpected signing method", fn), "current token", token.Header["alg"])
			return nil, ErrSigningMethod
		}

		return []byte(os.Getenv("gophermart_key")), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", ErrTokenExpired
		}

		j.log.Warn(fmt.Sprintf("%s: failed to parse token", fn), "error", err.Error())
		return "", ErrParseToken
	}

	if !token.Valid {
		j.log.Warn(fmt.Sprintf("%s: token is invalid", fn))
		return "", ErrInvalidToken
	}

	return claims.UserID, nil
}
