package main

import (
  "context"
  "crypto/sha256"
  "fmt"
  "log"
  "os"

  "github.com/aws/aws-lambda-go/events"
  "github.com/aws/aws-lambda-go/lambda"

  "github.com/aws/aws-sdk-go-v2/aws"
  "github.com/aws/aws-sdk-go-v2/config"
  "github.com/aws/aws-sdk-go-v2/credentials"
  "github.com/aws/aws-sdk-go-v2/service/dynamodb"
  "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
  "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Generic lambda response
type Response struct {
  Message string `json:"message"`
}

type Birthday struct {
  Id   string `dynamodbav:"Id"`
}

func (b Birthday) GetKey() map[string]types.AttributeValue {
  id, err := attributevalue.Marshal(b.Id)
  if err != nil {
    panic(err)
  }
  return map[string]types.AttributeValue{"Id": id}
}

var (
  tableName string
  db        *dynamodb.Client
)

func CreateLocalClient() *dynamodb.Client {
  cfg, err := config.LoadDefaultConfig(context.TODO(),
              config.WithRegion("ddblocal"),
              config.WithEndpointResolver(aws.EndpointResolverFunc(
                func(service, region string) (aws.Endpoint, error) {
                  return aws.Endpoint{URL: "http://dynamo:8000"}, nil
                })),
              config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
                Value: aws.Credentials{
                  AccessKeyID: "dummy",
                  SecretAccessKey: "dummy",
                  SessionToken: "dummy",
                  Source: "Local DynamoDB",
                },
              }))
  if err != nil {
    panic(err)
  }

  return dynamodb.NewFromConfig(cfg)
}

func init() {
  if "AWS_SAM_LOCAL" == os.Getenv("AWSENV") {
    log.Println("using local config")
    db = CreateLocalClient()
  } else {
    sdkConfig, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
      log.Fatal(err)
    }

    db = dynamodb.NewFromConfig(sdkConfig)
  }

  tableName = os.Getenv("DYNAMODB_TABLE")
  log.Println("Table is", tableName)
}

func Handler(request events.APIGatewayProxyRequest) (*Response, error) {
  user := request.PathParameters["username"]
  hash := sha256.Sum256([]byte(user))
  hashStr := fmt.Sprintf("%x", hash[:])

  fmt.Println(hashStr)

  bday := Birthday{Id: hashStr}

  result, err := db.GetItem(context.TODO(), &dynamodb.GetItemInput{
    TableName: aws.String(tableName),
    Key: bday.GetKey(),
  })
  if err != nil {
    msg := fmt.Errorf("GetItem: %s. error: %s", hashStr, err.Error())
    return nil, msg
  }
  if result.Item == nil {
    msg := fmt.Errorf("notFound: %s", user)
    return nil, msg
  }

  fmt.Printf("Dynamo: %#v\n", result.Item)

  return &Response{
    Message: fmt.Sprintf("Hello from Go, %s!", user),
  }, nil
}

func main() {
  lambda.Start(Handler)
}
