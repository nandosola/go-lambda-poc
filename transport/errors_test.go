package transport

import (
  "errors"
  "fmt"
  "testing"

  "service"

  "github.com/aws/aws-lambda-go/events"
  "github.com/kinbiko/jsonassert"
)

type mockedResponse struct {
  body   string
  cType  string
  status int
}

func TestErrorResponse(t *testing.T) {
  ja := jsonassert.New(t)
  mockReq := &Request{
    RequestContext: events.APIGatewayV2HTTPRequestContext{ RequestID: "12345-abcd" },
  }

  cases := []struct {
    name     string
    err      error
    expected mockedResponse
  }{
    {
      name: "ErrPathParamNotFound",
      err: fmt.Errorf("TestFoo: %s, %w", "fooparam", ErrPathParamNotFound),
      expected: mockedResponse{
        body: `{"message": "missing path param", "requestId": "12345-abcd"}`,
        cType: "application/json",
        status: 503,
      },
    },
    {
      name: "ValidationErrorRead",
      err: fmt.Errorf("ValidationError.Read: %w", errors.New("guru meditation")),
      expected: mockedResponse{
        body: `{"message": "username must be alphanumeric, min=2, max=12 chars", "requestId": "12345-abcd"}`,
        cType: "application/json",
        status: 400,
      },
    },
    {
      name: "ValidationErrorUpdate",
      err: fmt.Errorf("ValidationError.Update: %w", errors.New("This is so wrong")),
      expected: mockedResponse{
        body: `{"message": "date must be yyyy-mm-dd", "requestId": "12345-abcd"}`,
        cType: "application/json",
        status: 400,
      },
    },
    {
      name: "ServiceErrNotFound",
      err: fmt.Errorf("TestBar: %s, %w", "alice", service.ErrNotFound),
      expected: mockedResponse{
        body: `{"message": "username not found", "requestId": "12345-abcd"}`,
        cType: "application/json",
        status: 404,
      },
    },
    {
      name: "ServiceErrInvalidBirthday",
      err: fmt.Errorf("TestBaz: %w", service.ErrInvalidBirthday),
      expected: mockedResponse{
        body: `{"message": "birth date must be before today", "requestId": "12345-abcd"}`,
        cType: "application/json",
        status: 400,
      },
    },
    {
      name: "JSONUnmarshal",
      err: fmt.Errorf("JSONUnmarshal: %w", errors.New("I hate yaml")),
      expected: mockedResponse{
        body: `{"message": "bad json input", "requestId": "12345-abcd"}`,
        cType: "application/json",
        status: 400,
      },
    },
    {
      name: "GenericDBError",
      err: fmt.Errorf("DDBFooBarBaz: fatal crash"),
      expected: mockedResponse{
        body: `{"message": "database error", "requestId": "12345-abcd"}`,
        cType: "application/json",
        status: 500,
      },
    },
    {
      name: "DDBTimeout",
      err: fmt.Errorf("DDBContext: canceled"),
      expected: mockedResponse{
        body: `{"message": "database error", "requestId": "12345-abcd"}`,
        cType: "application/json",
        status: 503,
      },
    },
    {
      name: "UnexectedError",
      err: fmt.Errorf("BoomCrashBang: this is fine"),
      expected: mockedResponse{
        body: `{"message": "internal error", "requestId": "12345-abcd"}`,
        cType: "application/json",
        status: 500,
      },
    },
  }

  for _, tc := range cases {
    t.Run(tc.name, func(t *testing.T) {
      res, _ := ErrorResponse(tc.err, mockReq)
      if res.StatusCode != tc.expected.status {
        t.Errorf("expected status=%d, got %d", tc.expected.status, res.StatusCode)
      }
      if v, ok := res.Headers["Content-Type"]; !ok || v != tc.expected.cType {
        t.Errorf("expected content-type=%s, got '%s'", tc.expected.cType, v)
      }

      ja.Assertf(string(res.Body), tc.expected.body)
    })
  }
}

