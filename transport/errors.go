package transport

import (
  "log"
  "net/http"
  "strings"
)


const (
  infoLogLevel  = "INFO"
  warnLogLevel  = "WARN"
  errorLogLevel = "ERROR"
)

// Common eerror router
func ErrorResponse(err error, req *Request) (*Response, error) {
  reqId := req.RequestContext.RequestID
  errMsg := err.Error()

  userMsg := "internal error"
  status := http.StatusInternalServerError
  level := errorLogLevel

  switch {
  case strings.HasPrefix(errMsg, "DDBNotFound"):
    userMsg = "username not found"
    status  = http.StatusNotFound
    level   = infoLogLevel

  case strings.Contains(errMsg, "Field validation for 'Name'"):
    userMsg = "username must be alphanumeric, min=2, max=12 chars"
    status  = http.StatusBadRequest
    level   = warnLogLevel

  case strings.HasPrefix(errMsg, "PathParamNotFound"):
    userMsg = "username path param is missing"
    status  = http.StatusServiceUnavailable

  case strings.HasPrefix(errMsg, "DDB"):
    userMsg = "database error"
  }

  httpError :=  HttpError{
    Status:    status,
    Wrapped:   err,
    Message:   userMsg,
    RequestId: reqId,
  }

  log.Printf("%s - %s", level, httpError.Error())

  return httpError.asResponse(req)
}

