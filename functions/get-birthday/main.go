package main

import (
  "context"
  "log"

  "service"
  "transport"

  "github.com/aws/aws-lambda-go/lambda"
)


func init() {
  if err := service.InitializeRepo(); err != nil {
    log.Fatalf("the repository could not be initialized, %s", err.Error())
  }
}

func processGet(ctx context.Context, req transport.Request) (*transport.Response, error) {
  //log.Printf("Received req %#v", req)

  rr, err := transport.NewReadRequest(ctx, &req)
  if err != nil {
    return transport.ErrorResponse(err, &req)
  }

  res, err := service.Reader().GetBirthday(ctx, rr)
  if err != nil {
    return transport.ErrorResponse(err, &req)
  }

  return transport.SuccessResponse(res, &req)
}

func main() {
  lambda.Start(processGet)
}

