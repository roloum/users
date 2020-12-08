package aws

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

//GetSession returns an AWS session
func GetSession(log *log.Logger) (*session.Session, error) {

	log.Printf("Retrieving AWS Session")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
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
