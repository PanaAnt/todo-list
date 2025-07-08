package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secret = []byte("secret")

func GenerateJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyJWT(tokenString string) (*jwt.Token, error) {
	tokenn, err := jwt.Parse(tokenString, func(tokenn *jwt.Token) (interface{}, error) {
		if _, ok := tokenn.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing error: %v", tokenn.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	return tokenn, nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authheader := r.Header.Get("Authorization")
		if authheader == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authheader, "Bearer ")

		token, err := VerifyJWT(tokenString)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "invalid token claims", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["userID"].(string)
		if !ok {
			http.Error(w, "invalid user ID in token", http.StatusUnauthorized)
			return
		}

		context := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(context))
	})
}
