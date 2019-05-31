package app

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"strings"
	"two-server/models"
	u "two-server/utils"
)

var JwtAuthentication = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/user/register" || r.URL.Path == "/api/user/login" {
			next.ServeHTTP(w, r)
		} else {
			split := strings.Split(r.Header.Get("Authorization"), " ")
			if len(split) == 2 {
				tk := &models.Token{}
				token, _ := jwt.ParseWithClaims(split[1], tk, func(token *jwt.Token) (interface{}, error) {
					return []byte(os.Getenv("app_secret")), nil
				})
				if token.Valid {
					ctx := context.WithValue(r.Context(), "user", tk.UserId)
					r = r.WithContext(ctx)
					next.ServeHTTP(w, r)
					return
				}
			}
			u.Respond(w, 403, nil)
		}
	})
}