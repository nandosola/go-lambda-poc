package main

import (
  "fmt"
  "os"

  "github.com/aws/aws-lambda-go/events"
  "github.com/aws/aws-lambda-go/lambda"
)

// Generic lambda response
type Response struct {
  Message string `json:"message"`
}

var table string

func init() {
  table = os.Getenv("DYNAMODB_TABLE")
}

func Handler(request events.APIGatewayProxyRequest) (*Response, error) {
  user := request.PathParameters["username"]

  fmt.Println("Dynamo: ", table)

  return &Response{
    Message: fmt.Sprintf("Hello from Go, %s!", user),
  }, nil
}

func main() {
  lambda.Start(Handler)
}
