package middleware

import (
	"context"
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
)

type AuthMiddleware struct {
	AuthClient *auth.Client
}

func NewAuthMiddleware(client *auth.Client) *AuthMiddleware {
	return &AuthMiddleware{
		AuthClient: client,
	}
}

func (am *AuthMiddleware) Verify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// TODO tratar o erro aí

			return
		}

		splitToken := strings.Split(authHeader, " ")
		if len(splitToken) != 2 || strings.ToLower(splitToken[0]) != "bearer" {
			// TODO tratar o erro aí
			
			return
		}

		idToken := splitToken[1]

		token, err := am.AuthClient.VerifyIDToken(r.Context(), idToken)
		if err != nil {
			// TODO tratar o erro aí

			return
		}

		ctx := context.WithValue(r.Context(), "user", token)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}