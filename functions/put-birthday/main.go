package main

import (
  "context"

  "transport"

  "github.com/aws/aws-lambda-go/lambda"
)

func processPut(ctx context.Context, req transport.Request) (*transport.Response, error) {
 res := struct{
   Message string `json:"message"`
 }{
   Message: "not implemented yet",
 }

 return transport.SuccessResponse(res, &req)
}

func main() {
  lambda.Start(processPut)
}

