package middleware

import (
	"context"
	"go/adv-demo/configs"
	"go/adv-demo/pkg/jwt"
	"net/http"
	"strings"
)

type key string

const (
	ContextEmailKey key = "ContextEmailKey"
)

func writeUnauthed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
}

func IsAuthed(next http.Handler, conf *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authedToken := r.Header.Get("Authorization")
		if !strings.HasPrefix(authedToken, "Bearer") {
			writeUnauthed(w)
			return
		}

		token := strings.TrimPrefix(authedToken, "Bearer ")
		isValid, data := jwt.NewJwt(conf.Auth.Secret).Parse(token)

		if !isValid {
			writeUnauthed(w)
			return
		}

		ctx := context.WithValue(r.Context(), ContextEmailKey, data.Email)
		req := r.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}
