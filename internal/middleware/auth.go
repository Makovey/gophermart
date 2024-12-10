package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Makovey/gophermart/internal/logger"
	"github.com/Makovey/gophermart/pkg/jwt"
)

const (
	CtxUserIDKey = "UserIDKey"
)

type Auth struct {
	jwt *jwt.JWT
	log logger.Logger
}

func NewAuth(jwt *jwt.JWT, log logger.Logger) Auth {
	return Auth{jwt: jwt, log: log}
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

		userID, err := a.jwt.ParseUserID(authHeader)
		if err != nil {
			switch {
			case errors.Is(err, jwt.ErrParseToken):
				a.log.Info(fmt.Sprintf("%s: failed to parse token", fn), "token", authHeader)
				responseWithError(w, http.StatusInternalServerError, "internal server error, please try again")
				return
			case errors.Is(err, jwt.ErrSigningMethod),
				errors.Is(err, jwt.ErrInvalidToken),
				errors.Is(err, jwt.ErrTokenExpired):
				a.log.Info(fmt.Sprintf("%s: token is invalid", fn), "token", authHeader)
				responseWithError(w, http.StatusUnauthorized, "please, relogin again, to get access to this resource")
				return
			}
		}

		ctx := context.WithValue(r.Context(), CtxUserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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
