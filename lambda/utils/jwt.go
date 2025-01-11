package utils

import (
	"fmt"
	"lambda/types"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateToken(user *types.User) string {
	validUntil := time.Now().Add(time.Hour * 24).Unix()

	claims := jwt.MapClaims{
		"user":    user.Username,
		"expires": validUntil,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("JWT_SECRET")

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return ""
	}

	return tokenString
}

func ExtractTokenFromAuthHeader(
	headers map[string]string,
) string {
	authHeader, ok := headers["Authorization"]

	if !ok {
		return ""
	}

	splitToken := strings.Split(authHeader, "Bearer ")

	if len(splitToken) != 2 {
		return ""
	}

	return splitToken[1]
}

func ParseToken(tokenString string) (jwt.MapClaims, error) {
	secret := os.Getenv("JWT_SECRET")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
