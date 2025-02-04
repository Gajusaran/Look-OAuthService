package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Gajusaran/Look-OAuthService/util"
)

type ContextKey string

const ClaimsKey ContextKey = "claims" //is a custom type to avoid accidental collisions with other string keys. if we using built in string it's giving warning this is best practice

func TokenValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization") //access token
		if authHeader == "" {
			http.Error(w, "Authorization token missing", http.StatusUnauthorized)
			return
		}

		// Token is expected to be in the form "Bearer <token>"
		parts := strings.Split(authHeader, " ") //spilt the bearer and token
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid token format", http.StatusUnauthorized) // handling of format of token
			return
		}

		tokenStr := parts[1]

		claims, err := util.ParseToken(tokenStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized) // invalid token case
			return
		}

		// Attach the claims to the request context for further use in the handler
		ctx := context.WithValue(r.Context(), ClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
