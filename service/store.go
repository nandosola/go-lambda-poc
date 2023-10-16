package service

import (
  "context"
  "fmt"
  "os"
  "sync"
  "time"

  "github.com/aws/aws-sdk-go-v2/aws"
  "github.com/aws/aws-sdk-go-v2/config"
  "github.com/aws/aws-sdk-go-v2/credentials"
  "github.com/aws/aws-sdk-go-v2/service/dynamodb"
  "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
)



var (
  ddbStore     dynamoStore
  onceDdbStore sync.Once

  storeTimeoutSecs = 5 * time.Second
)

type dynamoStore struct {
  client    *dynamodb.Client
  tableName string
}

func storeContext() (context.Context, context.CancelFunc) {
  return context.WithTimeout(context.Background(), storeTimeoutSecs)
}

func createLocalClient() (*dynamodb.Client, error) {

  cfg, err := config.LoadDefaultConfig(context.Background(),  //non-cancellable
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
  var err error

  onceDdbStore.Do(func(){
    var db *dynamodb.Client

    if "AWS_SAM_LOCAL" == os.Getenv("AWSENV") {
      db, err = createLocalClient()
    } else {
      conf, errConf := config.LoadDefaultConfig(context.Background())
      if errConf != nil {
        err = fmt.Errorf("DDBConnect, AWSConfig: %w", err)
        return
      }
      db = dynamodb.NewFromConfig(conf)
    }

    ddbStore = dynamoStore{client: db, tableName: os.Getenv("DYNAMODB_TABLE")}
  })

  return &ddbStore, err
}

func (ds *dynamoStore) GetFromStore(bday *Birthday) (bool, error) {
  ctx, cancel := storeContext()
  defer cancel()

  result, err := ds.client.GetItem(ctx, &dynamodb.GetItemInput{
    TableName: aws.String(ds.tableName),
    Key: bday.GetKey(),
  })
  if err != nil {
    if ctx.Err() != nil {
      return false, fmt.Errorf("DDBContext: %w", err)
    }
    return false, fmt.Errorf("DDBGetItem: %s. Wrapped: %w", bday.Id, err)
  }
  if result.Item == nil {
    return false, nil
  }

  if err = attributevalue.UnmarshalMap(result.Item, bday); err != nil {
    return false, fmt.Errorf("DDBUnmarshalMap. Wrapped: %w", err)
  }

  return true, nil
}

func (ds *dynamoStore) AddToStore(bday *Birthday) error {
  item, err := attributevalue.MarshalMap(*bday)
  if err != nil {
    return fmt.Errorf("DDBMarshalMap: %s. Wrapped: %w", bday.Id, err)
  }

  ctx, cancel := storeContext()
  defer cancel()

  _, err = ds.client.PutItem(ctx, &dynamodb.PutItemInput{
    TableName: aws.String(ds.tableName),
    Item: item,
  })
  if err != nil {
    if ctx.Err() != nil {
      return fmt.Errorf("DDBContext: %w", err)
    }
    return fmt.Errorf("DDBPutItem: %s. Wrapped: %w", bday.Id, err)
  }

  return nil
}

