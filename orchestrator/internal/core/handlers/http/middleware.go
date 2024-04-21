package handlers

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

const (
	authorizationTokenCookie = "Authorization-Token"
	secretKey                = "secret-key"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

func WithAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		authToken, err := req.Cookie(authorizationTokenCookie)
		if err != nil {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		userID, err := getUserIDFromJWT(authToken.Value, secretKey)
		if err != nil {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		if userID == "" {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(req.Context(), "user_id", userID)
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}

func buildJWTString(userID, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func getUserIDFromJWT(tokenString, secretKey string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return "", err
	}

	return claims.UserID, nil
}
