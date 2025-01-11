package main

import (
	"fmt"
	"lambda/app"
	"lambda/middleware"
	"lambda/types"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(event types.Event) (string, error) {
	if event.Username == "" {
		return "", fmt.Errorf("username is required")
	}

	return fmt.Sprintf("Successfully called lambda function with username: %s", event.Username), nil
}

func ProtectedHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Protected route accessed",
	}, nil
}

func main() {
	myApp := app.NewApp()
	lambda.Start(
		func(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
			switch event.Path {
			case "/register":
				return myApp.ApiHandler.RegisterUserHandler(event)
			case "/login":
				return myApp.ApiHandler.LoginUser(event)
			case "/protected":
				return middleware.ValidateJwt(ProtectedHandler)(event)
			default:
				return events.APIGatewayProxyResponse{
					StatusCode: 404,
					Body:       "Not Found",
				}, nil
			}
		},
	)

}
