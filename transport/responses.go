package transport

import (
  "encoding/json"
  "fmt"
  "net/http"

  "github.com/aws/aws-lambda-go/events"
)

type (
  Request  events.APIGatewayV2HTTPRequest
  Response events.APIGatewayV2HTTPResponse
)

type HttpError struct {
  Status    int    `json:"-"`
  Wrapped   error  `json:"-"`
  Message   string `json:"message"`
  RequestId string `json:"requestId"`
}

func (err HttpError) Error() string {
  return fmt.Sprintf("status: %d, message: %s, requestId: %s, error: %s", err.Status, err.Message, err.RequestId, err.Wrapped.Error())
}

func (err HttpError) asResponse(req *Request) (*Response, error) {
  jsonStr, jsonErr := asJson(err)
  if jsonErr != nil {
    return ErrorResponse(jsonErr, req)
  }

  return &Response{
    StatusCode: err.Status,
    Body:       jsonStr,
    Headers: map[string]string{
      "Content-Type": "application/json",
    },
  }, nil
}

func SuccessResponse(res any, req *Request) (*Response, error) {
  jsonStr, err := asJson(res)
  if err != nil {
    return ErrorResponse(err, req)
  }

  return &Response{
    StatusCode: http.StatusOK,
    Body:       jsonStr,
    Headers: map[string]string{
      "Content-Type": "application/json",
    },
  }, nil
}

func asJson(obj any) (string, error) {
  jsonData, err  := json.Marshal(obj)
  if err != nil {
    return "", err
  }
  return string(jsonData), nil
}

