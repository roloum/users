package test

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

//Initializes Test environment
func init() {
	fmt.Println("Initializing test environment...")
}

//MockDynamoDB Mock DynamoDB client
type MockDynamoDB struct {
	dynamodbiface.DynamoDBAPI

	QueryOutput                *dynamodb.QueryOutput
	UpdateItemWithOutput       *dynamodb.UpdateItemOutput
	PutItemOutput              *dynamodb.PutItemOutput
	QueryOutputError           error
	UpdateItemWithContextError error
	PutItemOutputError         error
}

//PutItemWithContext mocks the PutItemWithContext function
func (m *MockDynamoDB) PutItemWithContext(aws.Context, *dynamodb.PutItemInput,
	...request.Option) (*dynamodb.PutItemOutput, error) {
	return m.PutItemOutput, m.PutItemOutputError
}
