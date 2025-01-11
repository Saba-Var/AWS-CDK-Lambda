package errors

import (
	"errors"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

var (
	ErrInvalidRequest     = errors.New("invalid request")
	ErrMissingCredentials = errors.New("username and password are required")
	ErrUserExists         = errors.New("user already exists")
)

var (
	ErrorResponse = map[error]events.APIGatewayProxyResponse{
		ErrInvalidRequest: {
			StatusCode: http.StatusBadRequest,
			Body:       "Invalid request",
		},
		ErrMissingCredentials: {
			StatusCode: http.StatusBadRequest,
			Body:       "Username and password are required",
		},
		ErrUserExists: {
			StatusCode: http.StatusConflict,
			Body:       "User already exists",
		},
	}

	InternalServerError = events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       "Internal server error",
	}

	SuccessResponse = events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       "User registered successfully",
	}

	InvalidResponse = events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Body:       "Invalid request",
	}

	UnauthorizedResponse = events.APIGatewayProxyResponse{
		StatusCode: http.StatusUnauthorized,
		Body:       "Unauthorized",
	}
)
