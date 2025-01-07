package main

import (
	"fmt"

	"lambda/app"

	"github.com/aws/aws-lambda-go/lambda"
)

type Event struct {
	Username string `json:"username"`
}

func HandleRequest(event Event) (string, error) {
	if event.Username == "" {
		return "", fmt.Errorf("username is required")
	}

	return fmt.Sprintf("Successfully called lambda function with username: %s", event.Username), nil
}

func main() {
	myApp := app.NewApp()
	lambda.Start(myApp.ApiHandler.RegisterUserHandler)
}
