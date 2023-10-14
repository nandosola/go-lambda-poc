package main

import (
  "context"

  "transport"

  "github.com/aws/aws-lambda-go/lambda"
)




func processPut(ctx context.Context, req transport.Request) (*transport.Response, error) {
 reqId := req.RequestContext.RequestID
 res := struct{
   Message string `json:"message"`
 }{
   Message: "not implemented yet",
 }

 return transport.SuccessResponse(res, reqId)
}

func main() {
  lambda.Start(processPut)
}

