package aws

import (
	"encoding/json"

	"github.com/rs/zerolog/log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

//GetSession returns an AWS session
func GetSession(region string) (*session.Session, error) {

	log.Debug().Msg("Retrieving AWS Session")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, err
	}

	return sess, nil

}

//GetDynamoDB returns an instance of the DynamoDB connection
func GetDynamoDB(sess *session.Session) *dynamodb.DynamoDB {
	dynamoSvc := dynamodb.New(sess)

	return dynamoSvc
}

// UnmarshalStreamImage converts events.DynamoDBAttributeValue to struct
func UnmarshalStreamImage(image map[string]events.DynamoDBAttributeValue,
	out interface{}) error {

	attributeMap := make(map[string]*dynamodb.AttributeValue)

	for k, v := range image {

		var dbAttr dynamodb.AttributeValue

		bytes, marshalErr := v.MarshalJSON()
		if marshalErr != nil {
			return marshalErr
		}

		json.Unmarshal(bytes, &dbAttr)
		attributeMap[k] = &dbAttr
	}

	return dynamodbattribute.UnmarshalMap(attributeMap, &out)

}
