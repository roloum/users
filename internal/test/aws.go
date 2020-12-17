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

	PutItemOutput            *dynamodb.PutItemOutput
	TransactWriteItemsOutput *dynamodb.TransactWriteItemsOutput
	OutputError              error
}

//PutItemWithContext mocks the PutItemWithContext method
func (m *MockDynamoDB) PutItemWithContext(aws.Context, *dynamodb.PutItemInput,
	...request.Option) (*dynamodb.PutItemOutput, error) {
	return m.PutItemOutput, m.OutputError
}

//TransactWriteItemsWithContext mocks the TransactWriteItemsWithContext method
func (m *MockDynamoDB) TransactWriteItemsWithContext(aws.Context,
	*dynamodb.TransactWriteItemsInput, ...request.Option) (
	*dynamodb.TransactWriteItemsOutput, error) {
	return m.TransactWriteItemsOutput, m.OutputError
}
