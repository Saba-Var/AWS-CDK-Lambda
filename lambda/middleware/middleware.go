package middleware

import (
	"fmt"
	"lambda/errors"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt"
)

type MiddlewareFunc = func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func ValidateJwt(
	next MiddlewareFunc,
) MiddlewareFunc {

	return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		token := extractTokenFromAuthHeader(request.Headers)

		if token == "" {
			return errors.UnauthorizedResponse, nil
		}

		claims, err := parseToken(token)

		if err != nil {
			return errors.UnauthorizedResponse, nil
		}

		expires := int64(claims["expires"].(float64))

		if time.Now().Unix() > expires {
			return errors.UnauthorizedResponse, nil
		}

		return next(request)
	}
}

func extractTokenFromAuthHeader(
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

func parseToken(tokenString string) (jwt.MapClaims, error) {
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
