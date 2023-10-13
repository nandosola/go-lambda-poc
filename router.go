package main

import (
  "context"
  "encoding/json"
  "log"
  "net/http"

  "github.com/aws/aws-lambda-go/events"
)


func router(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
  //log.Printf("Received req %#v", req)

  switch req.HTTPMethod {
  case "GET":
    return processGet(ctx, req)
  default:
    return clientError(http.StatusMethodNotAllowed)
  }
}

func processGet(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
  name, ok := req.PathParameters["username"]
  if !ok {
    return clientError(http.StatusNotFound)
  }

  rr, err := newReadRequest(ctx, name)
  if err != nil {
    log.Printf("Invalid request: %v", err)
    return clientError(http.StatusBadRequest)
  }

  res, err := GetBirthday(ctx, rr)
  if err != nil {
    return serverError(err)
  }

  return successResponse(res)
}

func clientError(status int) (*events.APIGatewayProxyResponse, error) {
  return &events.APIGatewayProxyResponse{
    Body:       http.StatusText(status),
    StatusCode: status,
  }, nil
}

func serverError(err error) (*events.APIGatewayProxyResponse, error) {
  log.Println(err.Error())

  return &events.APIGatewayProxyResponse{
    Body:       http.StatusText(http.StatusInternalServerError),
    StatusCode: http.StatusInternalServerError,
  }, nil
}

func successResponse(res *Birthday) (*events.APIGatewayProxyResponse, error) {
  json, err := json.Marshal(res)
  if err != nil {
    return serverError(err)
  }

  return &events.APIGatewayProxyResponse{
    StatusCode: http.StatusOK,
    Body:       string(json),
  }, nil
}

