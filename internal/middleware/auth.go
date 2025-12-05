package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/response"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("rahasia-negara-api")

type contextKey string
const UserIDKey contextKey = "userID"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.NewAPIError(http.StatusUnauthorized, "Token tidak ditemukan")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.NewAPIError(http.StatusUnauthorized, "Format token salah")
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			response.NewAPIError(http.StatusUnauthorized, "Token tidak valid atau kadaluwarsa")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.NewAPIError(http.StatusUnauthorized, "Token claims invalid")
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			response.NewAPIError(http.StatusUnauthorized, "User ID tidak ditemukan di token")
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next(w, r.WithContext(ctx))
	}
}

func GetUserIDFromContext(ctx context.Context) string {
	val := ctx.Value(UserIDKey)
	if val == nil {
		return ""
	}
	return val.(string)
}