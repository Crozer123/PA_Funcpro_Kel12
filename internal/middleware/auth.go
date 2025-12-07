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
		userID, err := validateToken(r)
		if err != nil {
			response.WriteJSON(w, http.StatusUnauthorized, "Unauthorized: "+err.Error(), nil)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next(w, r.WithContext(ctx))
	}
}

func AuthMiddlewareOptional(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			next(w, r)
			return
		}

		userID, err := validateToken(r)
		if err != nil {
			response.WriteJSON(w, http.StatusUnauthorized, "Token Invalid/Expired: "+err.Error(), nil)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next(w, r.WithContext(ctx))
	}
}

func validateToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", jwt.ErrTokenMalformed
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", jwt.ErrTokenMalformed
	}

	tokenString := parts[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", jwt.ErrTokenSignatureInvalid
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", jwt.ErrTokenInvalidClaims
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", jwt.ErrTokenInvalidClaims
	}

	return userID, nil
}

func GetUserIDFromContext(ctx context.Context) string {
	val := ctx.Value(UserIDKey)
	if val == nil {
		return ""
	}
	return val.(string)
}