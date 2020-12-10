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
