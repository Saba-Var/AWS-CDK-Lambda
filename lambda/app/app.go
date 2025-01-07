package app

import (
	"lambda/api"
	"lambda/database"
)

type App struct {
	ApiHandler api.ApiHandler
}

func NewApp() *App {
	dynamoDb := database.NewDynamoDbClient()

	return &App{
		ApiHandler: api.NewApiHandler(dynamoDb),
	}
}
