package api

import (
	"encoding/json"
	"fmt"
	"lambda/database"
	"lambda/types"
	"lambda/utils"
	"net/http"

	"lambda/errors"

	"github.com/aws/aws-lambda-go/events"
)

type ApiHandler struct {
	dbStore database.UserStore
}

func NewApiHandler(dbStore database.UserStore) *ApiHandler {
	return &ApiHandler{
		dbStore: dbStore,
	}
}

func (api *ApiHandler) RegisterUserHandler(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	registerUser, err := api.parseAndValidateRequest(event)
	if err != nil {
		response, exists := errors.ErrorResponse[err]
		if exists {
			return response, err
		}
		return errors.InternalServerError, err
	}

	if err := api.checkUserExists(registerUser.Username); err != nil {
		response, exists := errors.ErrorResponse[err]
		if exists {
			return response, nil
		}
		return errors.InternalServerError, err
	}

	if err := api.createAndSaveUser(registerUser); err != nil {
		return errors.InternalServerError, err
	}

	return errors.SuccessResponse, nil
}

func (api *ApiHandler) LoginUser(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var loginRequest LoginRequest

	err := json.Unmarshal([]byte(event.Body), &loginRequest)

	if err != nil {
		return errors.InvalidResponse, err
	}

	user, err := api.dbStore.GetUser(loginRequest.Username)

	if err != nil {
		return errors.InternalServerError, err
	}

	passwordsMatch := utils.ComparePasswordHash(loginRequest.Password, user.PasswordHash)

	if !passwordsMatch {
		return errors.UnauthorizedResponse, nil
	}

	token := utils.CreateToken(&user)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       fmt.Sprintf(`{"token": "%s"}`, token),
	}, nil
}

func (api *ApiHandler) parseAndValidateRequest(event events.APIGatewayProxyRequest) (*types.RegisterUser, error) {
	var registerUser types.RegisterUser
	if err := json.Unmarshal([]byte(event.Body), &registerUser); err != nil {
		return nil, errors.ErrInvalidRequest
	}

	if registerUser.Password == "" || registerUser.Username == "" {
		return nil, errors.ErrMissingCredentials
	}

	return &registerUser, nil
}

func (api *ApiHandler) checkUserExists(username string) error {
	exists, err := api.dbStore.DoesUserExists(username)
	if err != nil {
		return err
	}
	if exists {
		return errors.ErrUserExists
	}
	return nil
}

func (api *ApiHandler) createAndSaveUser(registerUser *types.RegisterUser) error {
	newUser, err := types.NewUser(registerUser)
	if err != nil {
		return err
	}

	return api.dbStore.RegisterUser(*newUser)
}
