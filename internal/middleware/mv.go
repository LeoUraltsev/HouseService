package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/LeoUraltsev/HouseService/internal/gen"
	"github.com/LeoUraltsev/HouseService/internal/jwt"
)

type Key string

const (
	UserTypeContextKey Key = "user_type"
	UserIDContextKey       = "user_id"
)

const Bearer = "Bearer"
const AuthHeader = "Authorization"

type Middleware struct {
	JWT *jwt.JWT
}

func (m *Middleware) AuthMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		_, ok := ctx.Value(gen.BearerAuthScopes).([]string)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}
		s := strings.Split(r.Header.Get(AuthHeader), " ")
		if len(s) < 2 || s[0] != Bearer {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token := s[1]

		claims, err := m.JWT.ParseToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userType := claims.UserType
		userID := claims.UserID

		ctx = context.WithValue(ctx, UserTypeContextKey, userType)
		ctx = context.WithValue(ctx, UserIDContextKey, userID)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
