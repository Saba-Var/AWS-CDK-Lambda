package utils

import (
	"github.com/golang-jwt/jwt"
	"lambda/types"
	"os"
	"time"
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
