package main

import (
  "context"
  "fmt"
  "log"
  "os"

  "github.com/aws/aws-sdk-go-v2/aws"
  "github.com/aws/aws-sdk-go-v2/config"
  "github.com/aws/aws-sdk-go-v2/credentials"
  "github.com/aws/aws-sdk-go-v2/service/dynamodb"
  "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
)

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

func GetFromStore(bday *Birthday) error {
  result, err := db.GetItem(context.TODO(), &dynamodb.GetItemInput{
    TableName: aws.String(tableName),
    Key: bday.GetKey(),
  })
  if err != nil {
    return fmt.Errorf("DDBGetItem: %s. Wrapped: %w", bday.Id, err)
  }
  if result.Item == nil {
    return fmt.Errorf("DDBNotFound: %s", bday.Id)
  }

  if err = attributevalue.UnmarshalMap(result.Item, &bday); err != nil {
    return fmt.Errorf("DDBUnmarshalMap. Wrapped: %w", err)
  }

  return nil
}

