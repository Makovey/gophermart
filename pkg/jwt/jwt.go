package jwt

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Makovey/gophermart/internal/logger"
)

const (
	tokenExp = time.Hour * 24
)

var (
	ErrParseToken    = errors.New("failed to parse token")
	ErrSigningMethod = errors.New("unexpected signing method")
	ErrInvalidToken  = errors.New("invalid token")
	ErrTokenExpired  = errors.New("token is expired")
)

type JWT struct {
	log logger.Logger
}

func NewJWT(logger logger.Logger) *JWT {
	return &JWT{log: logger}
}

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

func (j JWT) ParseUserID(tokenString string) (string, error) {
	fn := "jwt.ParseToken"

	var claims Claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			j.log.Warn(fmt.Sprintf("[%s] unexpected signing method", fn), "current token", token.Header["alg"])
			return nil, ErrSigningMethod
		}

		return []byte(os.Getenv("gophermart_key")), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", ErrTokenExpired
		}

		j.log.Warn(fmt.Sprintf("[%s] failed to parse token", fn), "error", err.Error())
		return "", ErrParseToken
	}

	if !token.Valid {
		j.log.Warn(fmt.Sprintf("[%s] token is invalid", fn))
		return "", ErrInvalidToken
	}

	return claims.UserID, nil
}

func (j JWT) BuildNewJWT(userID string) (string, error) {
	fn := "jwt.BuildNewJWT:"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("gophermart_key")))
	if err != nil {
		j.log.Warn(fmt.Sprintf("[%s] can't sign token", fn), "error", err.Error())
		return "", err
	}

	return tokenString, nil
}
