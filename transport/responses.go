package transport

import (
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "strings"

  "github.com/aws/aws-lambda-go/events"
)

type (
  Request  events.APIGatewayV2HTTPRequest
  Response events.APIGatewayV2HTTPResponse
)

const errorJsonTmpl = `{"message":"%s","requestId":"%s"}`


type ErrorResponse struct {
  Status  int    `json:"-"`
  Wrapped error  `json:"-"`
}

func (er ErrorResponse) Error() string {
  var msg string
  if er.Wrapped != nil {
    msg = er.Wrapped.Error()
  }

  return fmt.Sprintf("status: %d, message: %s", er.Status, msg)
}

func (er ErrorResponse) AsResponse(requestId string) (*Response, error) {
  res := &Response{ StatusCode: er.Status }

  if msg := er.logAndTranslate(); msg != "" {
    res.Body = fmt.Sprintf(errorJsonTmpl, msg, requestId)
  }

  res.Headers = map[string]string{
    "Content-Type": "application/json",
  }

  return res, nil
}


func (er ErrorResponse) logAndTranslate() string {
  if er.Wrapped == nil {
    return ""
  }

  errMsg := er.Wrapped.Error()
  log.Println(errMsg)

  switch {
  case strings.HasPrefix(errMsg, "DDBNotFound"):
    return "username not found"
  case strings.Contains(errMsg, "Field validation for 'Name'"):
    return "username must be alphanumeric, min=2, max=12 chars"
  case strings.HasPrefix(errMsg, "DDB"):
    return "database error"
  }

  return ""
}

func InternalError(err error) ErrorResponse {
  status := http.StatusInternalServerError
  errMsg := err.Error()
  switch {
    case strings.HasPrefix(errMsg, "DDBNotFound"):
      status = http.StatusNotFound
  }

  return ErrorResponse{Status: status, Wrapped: err}
}

var (
  MethodNotAllowedError = ErrorResponse{Status: http.StatusMethodNotAllowed}
  NotFoundError         = ErrorResponse{Status: http.StatusNotFound}
  BadRequestError       = func(err error) ErrorResponse { return ErrorResponse{Status: http.StatusBadRequest, Wrapped: err} }
)

func SuccessResponse(res any, requestId string) (*Response, error) {
  json, err  := json.Marshal(res)
  if err != nil {
    return InternalError(err).AsResponse(requestId)
  }

  return &Response{
    StatusCode: http.StatusOK,
    Body:       string(json),
    Headers: map[string]string{
      "Content-Type": "application/json",
    },
  }, nil
}

