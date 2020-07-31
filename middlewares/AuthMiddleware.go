package middlewares

import (
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	Tokens []string
}

func (amw *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		splitToken := strings.Split(token, "Bearer ")
		if len(splitToken) > 1 {
			token = splitToken[1]
		}

		found := false
		for k := range amw.Tokens {
			if amw.Tokens[k] == token {
				found = true
				break
			}
		}

		if found {
			next.ServeHTTP(w, r)
		} else {
			if r.Method == http.MethodPost {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			r.Header.Add("Authorization", "unauthorized")
			next.ServeHTTP(w, r)
		}
	})
}
