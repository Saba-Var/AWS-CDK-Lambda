package database

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	USER_TABLE = "userTable"
)

type UserStore interface {
	DoesUserExists(username string) (bool, error)
	RegisterUser(username, password string) error
}

type DynamoDbClient struct {
	databaseStore *dynamodb.DynamoDB
	UserStore
}

func (db *DynamoDbClient) DoesUserExists(username string) (bool, error) {
	result, err := db.databaseStore.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(USER_TABLE),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: &username,
			},
		},
	})

	if err != nil {
		return false, err
	}

	if result.Item == nil {
		return false, nil
	}

	return true, nil
}

func (db *DynamoDbClient) RegisterUser(username, password string) error {
	_, err := db.databaseStore.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(USER_TABLE),
		Item: map[string]*dynamodb.AttributeValue{
			"username": {
				S: &username,
			},
			"password": {
				S: &password,
			},
		},
	})

	if err != nil {
		return err
	}

	fmt.Printf("user %s registered\n", username)

	return nil
}

func NewDynamoDbClient() *DynamoDbClient {
	return &DynamoDbClient{
		databaseStore: dynamodb.New(
			session.Must(session.NewSession()),
		),
	}
}
