package transport

import (
  "errors"
  "fmt"
  "log"
  "net/http"
  "strings"

  "service"
)


const (
  infoLogLevel  = "INFO"
  warnLogLevel  = "WARN"
  errorLogLevel = "ERROR"
)

// Common error handler, this is a Go idiom
func ErrorResponse(err error, req *Request) (*Response, error) {
  reqId := req.RequestContext.RequestID
  errMsg := err.Error()

  userMsg := "internal error"
  status := http.StatusInternalServerError
  level := errorLogLevel

  switch {
    case errors.Is(err, service.ErrNotFound):
    userMsg = errors.Unwrap(err).Error()
    status  = http.StatusNotFound
    level   = infoLogLevel

  case errors.Is(err, ErrPathParamNotFound):
    userMsg = errors.Unwrap(err).Error()
    status  = http.StatusServiceUnavailable

  case strings.HasPrefix(errMsg, "ValidationError.Read"):
    userMsg = fmt.Sprintf("username must be %s chars", readReqRestrictions)
    status  = http.StatusBadRequest
    level   = warnLogLevel

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

