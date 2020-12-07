package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	uaws "github.com/roloum/users/internal/aws"
	"github.com/roloum/users/internal/config"
	"github.com/roloum/users/internal/user"
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
		}
	}
)

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, dynamoDB *dynamodb.DynamoDB,
	request events.APIGatewayProxyRequest,
	cfg configuration,
	log *log.Logger) (Response, error) {

	var body createRequest
	if err := json.Unmarshal([]byte(request.Body), &body); err != nil {
		return getResponse(http.StatusUnprocessableEntity, err.Error(), nil, log)
	}

	newUser := &user.NewUser{
		Email:     body.Email,
		FirstName: body.FirstName,
		LastName:  body.LastName,
	}

	u, err := user.Create(ctx, dynamoDB, newUser, "User", log)
	if err != nil {
		return getResponse(http.StatusUnprocessableEntity, err.Error(), nil, log)
	}

	return getResponse(http.StatusCreated, "User created", u, log)
}

// getResponse builds an API Gateway Response
func getResponse(statusCode int, message string, u *user.User, log *log.Logger) (
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

	log := log.New(os.Stdout, "Users: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	//Config holds the configuration for the application
	var cfg configuration
	err := config.Load(&cfg, log)
	if err != nil {
		return Response{}, err
	}

	fmt.Println(cfg)

	sess, err := uaws.GetSession(log)
	if err != nil {
		return Response{}, err
	}

	return Handler(ctx, uaws.GetDynamoDB(sess), request, cfg, log)

}

func main() {
	lambda.Start(initHandler)
}
