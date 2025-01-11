package database

import (
	"fmt"
	"lambda/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	USER_TABLE = "userTable"
)

type UserStore interface {
	DoesUserExists(username string) (bool, error)
	RegisterUser(user types.User) error
	GetUser(username string) (types.User, error)
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

func (db *DynamoDbClient) RegisterUser(user types.User) error {
	_, err := db.databaseStore.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(USER_TABLE),
		Item: map[string]*dynamodb.AttributeValue{
			"username": {
				S: &user.Username,
			},
			"password": {
				S: &user.PasswordHash,
			},
		},
	})

	if err != nil {
		return err
	}

	fmt.Printf("user %s registered\n", user.Username)

	return nil
}

func (db *DynamoDbClient) GetUser(username string) (types.User, error) {
	var user types.User

	result, err := db.databaseStore.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(USER_TABLE),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	})

	if err != nil {
		return user, err
	}

	if result.Item == nil {
		return user, fmt.Errorf("user %s not found", username)
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &user)

	if err != nil {
		return user, err
	}

	return user, nil
}

func NewDynamoDbClient() *DynamoDbClient {
	return &DynamoDbClient{
		databaseStore: dynamodb.New(
			session.Must(session.NewSession()),
		),
	}
}
