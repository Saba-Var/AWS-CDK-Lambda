package middleware

import (
	"lambda/errors"
	"lambda/utils"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

type MiddlewareFunc = func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func ValidateJwt(
	next MiddlewareFunc,
) MiddlewareFunc {

	return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		token := utils.ExtractTokenFromAuthHeader(request.Headers)

		if token == "" {
			return errors.UnauthorizedResponse, nil
		}

		claims, err := utils.ParseToken(token)

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
