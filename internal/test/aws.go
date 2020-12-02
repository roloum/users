package test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

//MockDynamoDB Mock DynamoDB client
type MockDynamoDB struct {
	dynamodbiface.DynamoDBAPI

	PutItemOutput      *dynamodb.PutItemOutput
	PutItemOutputError error
}

//PutItemWithContext mocks the PutItemWithContext function
func (m *MockDynamoDB) PutItemWithContext(aws.Context, *dynamodb.PutItemInput,
	...request.Option) (*dynamodb.PutItemOutput, error) {
	return m.PutItemOutput, m.PutItemOutputError
}
