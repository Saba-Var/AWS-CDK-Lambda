package api

import (
	"encoding/json"
	"lambda/database"
	"lambda/types"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type ApiHandler struct {
	dbStore database.UserStore
}

func NewApiHandler(dbStore database.UserStore) ApiHandler {
	return ApiHandler{
		dbStore: dbStore,
	}
}

func (api ApiHandler) RegisterUserHandler(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var registerUser types.RegisterUser

	err := json.Unmarshal([]byte(event.Body), &registerUser)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Invalid request",
		}, err
	}

	if registerUser.Password == "" || registerUser.Username == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Username and password are required",
		}, err
	}

	userExists, err := api.dbStore.DoesUserExists(registerUser.Username)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Internal server error",
		}, err
	}

	if userExists {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusConflict,
			Body:       "User already exists",
		}, nil
	}

	newUser, err := types.NewUser(&registerUser)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Internal server error",
		}, err
	}

	err = api.dbStore.RegisterUser(*newUser)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Internal server error",
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       "User registered successfully",
	}, nil
}
