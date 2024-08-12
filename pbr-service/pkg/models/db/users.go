package db

import (
	"errors"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type User struct {
	UserID        string   `dynamodbav:"user_id"`
	Email         string   `dynamodbav:"email"`
	ConnectionIDs []string `dynamodbav:"connection_ids"`
}

func FetchUserByID(userID string) (*User, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	input := &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("USER_TABLE_ARN")),
		Key: map[string]*dynamodb.AttributeValue{
			"user_id": {
				S: aws.String(userID),
			},
		},
	}

	result, err := svc.GetItem(input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		log.Println("No user found for the specified user_id")
		return nil, errors.New("no user found for the specified user_id")
	}

	user := &User{}
	err = dynamodbattribute.UnmarshalMap(result.Item, user)
	log.Println(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
