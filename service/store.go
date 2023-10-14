package service

import (
  "context"
  "fmt"
  "os"

  "github.com/aws/aws-sdk-go-v2/aws"
  "github.com/aws/aws-sdk-go-v2/config"
  "github.com/aws/aws-sdk-go-v2/credentials"
  "github.com/aws/aws-sdk-go-v2/service/dynamodb"
  "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
)


var ddbStore dynamoStore

type dynamoStore struct {
  client    *dynamodb.Client
  tableName string
}

func createLocalClient() (*dynamodb.Client, error) {
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
    return nil, fmt.Errorf("DDBLocal: %w", err)
  }

  return dynamodb.NewFromConfig(cfg), nil
}

func ddbConnect() (*dynamoStore, error) {
  var db *dynamodb.Client

  if "AWS_SAM_LOCAL" == os.Getenv("AWSENV") {
    var err error
    db, err = createLocalClient()
    if err != nil {
      return nil, err
    }
  } else {
    conf, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
      return nil, fmt.Errorf("DDBConnect, AWSConfig: %w", err)
    }

    db = dynamodb.NewFromConfig(conf)
  }

  ddbStore = dynamoStore{client: db, tableName: os.Getenv("DYNAMODB_TABLE")}
  return &ddbStore, nil
}

func (ds *dynamoStore) GetFromStore(bday *Birthday) error {
  result, err := ds.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
    TableName: aws.String(ds.tableName),
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

