package jwt

import (
	"fmt"
	"time"

	gojwt "github.com/dgrijalva/jwt-go"
)

var secret = []byte("your-secret-key")

type Claims struct {
	UserID string `json:"user_id"`
	Login  string `json:"login"`
	gojwt.StandardClaims
}

func GenerateToken(userID, login string) (string, error) {
	claims := Claims{
		UserID: userID,
		Login:  login,
		StandardClaims: gojwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := gojwt.ParseWithClaims(tokenString, &Claims{}, func(token *gojwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*gojwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}