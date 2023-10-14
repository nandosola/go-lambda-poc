package main

import (
  "context"
  "log"

  "service"

  "github.com/aws/aws-lambda-go/events"
  "github.com/aws/aws-lambda-go/lambda"
)


type (
  LambdaRequest  events.APIGatewayV2HTTPRequest
  LambdaResponse events.APIGatewayV2HTTPResponse
)


func init() {
  if err := service.InitializeRepo(); err != nil {
    log.Fatalf("the repository could not be initialized, %s", err.Error())
  }
}

func processGet(ctx context.Context, req LambdaRequest) (*LambdaResponse, error) {
  //log.Printf("Received req %#v", req)
  reqId := req.RequestContext.RequestID

  name, ok := req.PathParameters["username"]
  if !ok {
    return notFoundError.asLambdaResponse(reqId)
  }

  rr, err := newReadRequest(ctx, name)
  if err != nil {
    return badRequestError(err).asLambdaResponse(reqId)
  }

  res, err := service.Reader().GetBirthday(ctx, rr)
  if err != nil {
    return internalError(err).asLambdaResponse(reqId)
  }

  return successResponse(res, reqId)
}

func main() {
  lambda.Start(processGet)
}

