package api

import (
	"fmt"
	"lambda/database"
	"lambda/types"
)

type ApiHandler struct {
	dbStore database.DynamoDbClient
}

func NewApiHandler(dbStore *database.DynamoDbClient) ApiHandler {
	return ApiHandler{
		dbStore: *dbStore,
	}
}

func (api ApiHandler) RegisterUserHandler(event *types.RegisterUser) error {
	if event.Password == "" || event.Username == "" {
		return fmt.Errorf("username and password are required")
	}

	userExists, err := api.dbStore.DoesUserExists(event.Username)

	if err != nil {
		return fmt.Errorf("error checking if user exists: %v", err)
	}

	if userExists {
		return fmt.Errorf("user already exists")
	}

	err = api.dbStore.RegisterUser(event.Username, event.Password)

	if err != nil {
		return fmt.Errorf("error registering user: %v", err)
	}

	return nil
}
