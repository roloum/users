//Lambda function that creates an user in dynamoDB
package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	uaws "github.com/roloum/users/internal/aws"
	"github.com/roloum/users/internal/config"
	"github.com/roloum/users/internal/user"
)

const (
	//MsgUserCreated message returned when user is created successfully
	MsgUserCreated = "UserCreated"
)

type (
	// createRequest
	createRequest struct {
		Email     string `json:"email,omitempty"`
		FirstName string `json:"firstName,omitempty"`
		LastName  string `json:"lastName,omitempty"`
	}

	// createResponse
	createResponse struct {
		StatusCode int        `json:"status"`
		Message    string     `json:"message"`
		User       *user.User `json:"user,omitempty"`
	}

	// Response is of type APIGatewayProxyResponse since we're leveraging the
	// AWS Lambda Proxy Request functionality (default behavior)
	//
	// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
	Response events.APIGatewayProxyResponse

	configuration struct {
		AWS struct {
			DynamoDB struct {
				Table struct {
					User string `required:"true"`
				}
			}
			Region string `required:"true"`
		}
	}
)

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, dynamoDB *dynamodb.DynamoDB,
	request events.APIGatewayProxyRequest,
	cfg configuration) (Response, error) {

	var body createRequest
	if err := json.Unmarshal([]byte(request.Body), &body); err != nil {
		return getResponse(http.StatusUnprocessableEntity, err.Error(), nil)
	}

	newUser := &user.NewUser{
		Email:     body.Email,
		FirstName: body.FirstName,
		LastName:  body.LastName,
	}

	u, err := user.Create(ctx, dynamoDB, newUser, "User")
	if err != nil {
		return getResponse(http.StatusUnprocessableEntity, err.Error(), nil)
	}

	return getResponse(http.StatusCreated, MsgUserCreated, u)
}

// getResponse builds an API Gateway Response
func getResponse(statusCode int, message string, u *user.User) (
	Response, error) {

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	resp := &createResponse{
		StatusCode: statusCode,
		Message:    message,
		User:       u,
	}

	js, err := json.Marshal(resp)
	if err != nil {
		return Response{StatusCode: http.StatusInternalServerError}, err
	}

	log.Printf("status_code: %d, message: %s", resp.StatusCode, resp.Message)

	return Response{Headers: headers, Body: string(js),
		StatusCode: resp.StatusCode}, nil
}

func initHandler(ctx context.Context, request events.APIGatewayProxyRequest) (
	Response, error) {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	//Config holds the configuration for the application
	var cfg configuration
	err := config.Load(&cfg)
	if err != nil {
		return Response{}, err
	}

	sess, err := uaws.GetSession(cfg.AWS.Region)
	if err != nil {
		return Response{}, err
	}

	return Handler(ctx, uaws.GetDynamoDB(sess), request, cfg)

}

func main() {
	lambda.Start(initHandler)
}
