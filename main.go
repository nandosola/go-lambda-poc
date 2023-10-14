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
  reqId := req.RequestContext.RequestID

  name, ok := req.PathParameters["username"]
  if !ok {
    return transport.NotFoundError.AsResponse(reqId)
  }

  rr, err := newReadRequest(ctx, name)
  if err != nil {
    return transport.BadRequestError(err).AsResponse(reqId)
  }

  res, err := service.Reader().GetBirthday(ctx, rr)
  if err != nil {
    return transport.InternalError(err).AsResponse(reqId)
  }

  return transport.SuccessResponse(res, reqId)
}

func main() {
  lambda.Start(processGet)
}

