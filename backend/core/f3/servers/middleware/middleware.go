package middleware

import (
	"context"
	"f3/models"
	"github.com/go-chi/render"
	jwt "github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

type UserAccessCtx struct {
	jwtKey []byte
}

func NewUserAccessCtx(jwtKey []byte) *UserAccessCtx {
	return &UserAccessCtx{
		jwtKey: jwtKey,
	}
}

// ChiMiddleware Middleware аунтификации
func (u UserAccessCtx) ChiMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtToken := r.Header.Get("Authorization")
		splitToken := strings.Split(jwtToken, "Bearer ")
		if len(splitToken) < 2 {
			_ = render.Render(w, r, models.Unauthorized(models.ErrTokenNotFound))
			return
		}

		jwtToken = splitToken[1]

		token, err := jwt.ParseWithClaims(jwtToken, &models.User{}, func(token *jwt.Token) (interface{}, error) {
			return u.jwtKey, nil
		})

		ctx := r.Context()

		if err != nil {
			_ = render.Render(w, r, models.Unauthorized(err))
			return
		} else if claims, ok := token.Claims.(*models.User); ok {
			ctx = context.WithValue(ctx, "name", claims.Name)
			ctx = context.WithValue(ctx, "email", claims.Email)
			ctx = context.WithValue(ctx, "tdid", claims.Tdid)
			ctx = context.WithValue(ctx, "login", claims.Login)
		} else {
			_ = render.Render(w, r, models.Unauthorized(err))
			return
		}

		// токен валидный, пропускаем его
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
