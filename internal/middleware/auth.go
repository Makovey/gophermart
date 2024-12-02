package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Makovey/gophermart/internal/logger"
)

var (
	errParseToken    = errors.New("failed to parse token")
	errSigningMethod = errors.New("unexpected signing method")
	errInvalidToken  = errors.New("invalid token")
	errTokenExpired  = errors.New("token is expired")
)

type Auth struct {
	log logger.Logger
}

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

func (a Auth) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fn := "auth.Authenticate"

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			a.log.Info(fmt.Sprintf("%s: missing authorization header", fn))
			responseWithError(w, http.StatusUnauthorized, "please, login to access this resource")
			return
		}

		userID, err := a.parseUserID(authHeader)
		if err != nil {
			switch {
			case errors.Is(err, errParseToken):
				a.log.Info(fmt.Sprintf("%s: failed to parse token", fn), "token", authHeader)
				responseWithError(w, http.StatusInternalServerError, "internal server error, please try again")
			case errors.Is(err, errSigningMethod),
				errors.Is(err, errInvalidToken),
				errors.Is(err, errTokenExpired):
				a.log.Info(fmt.Sprintf("%s: token is invalid", fn), "token", authHeader)
				responseWithError(w, http.StatusUnauthorized, "please, relogin again, to get access to this resource")
				return
			}
		}

		ctx := context.WithValue(r.Context(), "CtxUserIDKey", userID) // TODO: change to const
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a Auth) parseUserID(tokenString string) (string, error) {
	fn := "auth.parseToken"

	var claims Claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			a.log.Warn(fmt.Sprintf("%s: unexpected signing method", fn), "current token", token.Header["alg"])
			return nil, errSigningMethod
		}

		return []byte(os.Getenv("gophermart_key")), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", errTokenExpired
		}

		a.log.Warn(fmt.Sprintf("%s: failed to parse token", fn), "error", err.Error())
		return "", errParseToken
	}

	if !token.Valid {
		a.log.Warn(fmt.Sprintf("%s: token is invalid", fn))
		return "", errInvalidToken
	}

	return claims.UserID, nil
}

func responseWithError(w http.ResponseWriter, status int, message string) {
	type Response struct {
		Error string `json:"error"`
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := Response{
		Error: message,
	}

	json.NewEncoder(w).Encode(response)
}

func NewAuth(log logger.Logger) Auth {
	return Auth{log: log}
}
