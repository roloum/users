package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	uaws "github.com/roloum/users/internal/aws"
	"github.com/roloum/users/internal/config"
	"github.com/roloum/users/internal/user"
)

const (
	//MsgUserActivated message returned when user is activated successfully
	MsgUserActivated = "MsgUserActivated"

	//ErrorEmailIsEmpty message returned if email is empty
	ErrorEmailIsEmpty = "EmailIsEmpty"

	//ErrorTokenIsEmpty message returned if token is empty
	ErrorTokenIsEmpty = "TokenIsEmpty"
)

type (

	// activateResponse
	activateResponse struct {
		StatusCode int    `json:"status"`
		Message    string `json:"message"`
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

	email := request.QueryStringParameters["email"]
	if email == "" {
		return getResponse(http.StatusUnprocessableEntity, ErrorEmailIsEmpty)
	}

	token := request.QueryStringParameters["token"]
	if token == "" {
		return getResponse(http.StatusUnprocessableEntity, ErrorTokenIsEmpty)
	}

	email = strings.ToLower(email)
	log.Info().Msgf("Activating account: %s", email)

	u := &user.User{
		Email: email,
	}

	err := u.Activate(ctx, dynamoDB, cfg.AWS.DynamoDB.Table.User, token)
	if err != nil {
		return getResponse(http.StatusUnprocessableEntity, err.Error())
	}

	log.Info().Msg("User Activated")

	return getResponse(http.StatusCreated, MsgUserActivated)
}

// getResponse builds an API Gateway Response
func getResponse(statusCode int, message string) (Response, error) {

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	resp := &activateResponse{
		StatusCode: statusCode,
		Message:    message,
	}

	js, err := json.Marshal(resp)
	if err != nil {
		return Response{StatusCode: http.StatusInternalServerError}, err
	}

	log.Debug().Msgf("status_code: %d, message: %s", resp.StatusCode, resp.Message)

	return Response{Headers: headers, Body: string(js),
		StatusCode: resp.StatusCode}, nil
}

func initHandler(ctx context.Context, request events.APIGatewayProxyRequest) (
	Response, error) {

	//Config holds the configuration for the application
	var cfg configuration
	err := config.Load(&cfg)
	if err != nil {
		return Response{}, err
	}

	log.Debug().Msg("initHandler function")

	sess, err := uaws.GetSession(cfg.AWS.Region)
	if err != nil {
		return Response{}, err
	}

	return Handler(ctx, uaws.GetDynamoDB(sess), request, cfg)

}

func main() {
	lambda.Start(initHandler)
}
