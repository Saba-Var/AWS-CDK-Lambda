package database

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DynamoDbClient struct {
	databaseStore *dynamodb.DynamoDB
}

func NewDynamoDbClient() *DynamoDbClient {
	return &DynamoDbClient{
		databaseStore: dynamodb.New(
			session.Must(session.NewSession()),
		),
	}
}
