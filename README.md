DynamoDB project
This is a small serverless project that creates an user in a DynamoDB table. Via
DynamoDB streams, an email will be sent to the user in order to activate the
account.

Lambda functions:
 - createUser
 - notifyUser (triggered by DynamoDB stream)
 - activateUser

DynamoDB tables:
 - User

Serverless example
 - https://github.com/serverless/examples/blob/master/aws-golang-dynamo-stream-to-elasticsearch/serverless.yml
 - https://www.serverless.com/framework/docs/providers/aws/events/streams/
