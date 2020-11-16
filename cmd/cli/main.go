package main

import (
	"context"
	"log"

	"os"

	"github.com/roloum/users/cmd/cli/internal/cmd"
	uaws "github.com/roloum/users/internal/aws"
)

func main() {
	log := log.New(os.Stdout, "Users: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	if err := run(log); err != nil {
		log.Printf("Main: %s", err)
		os.Exit(1)
	}
}

func run(log *log.Logger) error {

	ctx := context.WithValue(context.Background(), cmd.ContextKey(cmd.LOG), log)

	sess, err := uaws.GetSession(log)
	if err != nil {
		return err
	}
	dynamo := uaws.GetDynamoDB(sess)
	ctx = context.WithValue(ctx, cmd.ContextKey(cmd.DYNAMO), dynamo)

	if err := cmd.RootCmd.ExecuteContext(ctx); err != nil {
		return err
	}

	/*
		u, err := user.Create(ctx, &user.NewUser{})

		if err != nil {
			return err
		}

		_ = u
	*/
	return nil
}

// 	fmt.Println("Hello again")
//
// 	sess, err := session.NewSession(&aws.Config{
// 		Region: aws.String("us-west-2"),
// 	})
// 	if err != nil {
// 		return err
// 	}
//
// 	svc := dynamodb.New(sess)
//
// 	user := User{
// 		"rolando.umana@gmail.com",
// 		time.Now().Format("2006-01-02"),
// 		Attributes{
// 			"Rolando",
// 			"Umana",
// 			*aws.Bool(false)}}
//
// 	av, err := dynamodbattribute.MarshalMap(user)
// 	if err != nil {
// 		return err
// 	}
//
// 	input := &dynamodb.PutItemInput{
// 		Item:      av,
// 		TableName: aws.String("User"),
// 	}
//
// 	_, err = svc.PutItem(input)
// 	if err != nil {
// 		return err
// 	}
//
// 	fmt.Println("Created...")
//
// 	return nil
// }
//
// //User ...
// type User struct {
// 	Email      string     `json:"email,omitempty"`
// 	Created    string     `json:"created,omitempty"`
// 	Attributes Attributes `json:"attributes,omitempty"`
// }
//
// //Attributes ...
// type Attributes struct {
// 	FirstName string `json:"firstName,omitempty"`
// 	LastName  string `json:"lastName,omitempty"`
// 	Active    bool   `type:"BOOL" json:"active,omitempty"`
// }
