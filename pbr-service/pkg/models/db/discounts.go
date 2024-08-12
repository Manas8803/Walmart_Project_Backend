package db

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DiscountDB struct {
	BeaconID      string `dynamodbav:"beacon_id"`
	DiscountOffer string `dynamodbav:"discount_offer"`
}

func FetchDiscountByBeaconID(beaconID string) (*DiscountDB, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	input := &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("DISCOUNT_TABLE_ARN")),
		Key: map[string]*dynamodb.AttributeValue{
			"beacon_id": {
				S: aws.String(beaconID),
			},
		},
	}

	result, err := svc.GetItem(input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, nil
	}

	discount := &DiscountDB{}
	err = dynamodbattribute.UnmarshalMap(result.Item, discount)
	if err != nil {
		return nil, err
	}

	return discount, nil
}
