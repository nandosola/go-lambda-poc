package main

import (
  "context"
  "log"

  "transport"
  "service"

  "github.com/aws/aws-lambda-go/lambda"
)

func init() {
  if err := service.InitializeRepo(); err != nil {
    log.Fatalf("the repository could not be initialized, %s", err.Error())
  }
}

func processPut(ctx context.Context, req transport.Request) (*transport.Response, error) {
  bday, err := transport.NewAddOrUpdateRequest(ctx, &req)
  if err != nil {
    return transport.ErrorResponse(err, &req)
  }

  if err := service.Writer().AddBirthday(ctx, bday); err != nil {
    return transport.ErrorResponse(err, &req)
  }

  return transport.NoContent()
}

func main() {
  lambda.Start(processPut)
}

