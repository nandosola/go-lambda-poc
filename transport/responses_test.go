package transport

import (
  "testing"

  "github.com/aws/aws-lambda-go/events"
  "github.com/kinbiko/jsonassert"
)

 
func TestSuccessResponse(t *testing.T) {
  ja := jsonassert.New(t)
  mockReq := &Request{
    RequestContext: events.APIGatewayV2HTTPRequestContext{ RequestID: "12345-abcd" },
  }

  wadus := struct{
    Thing  string  `json:"thing"`
    Answer int     `json:"answer"`
  }{
    Thing: "thingamabob",
    Answer: 42,
  }

  expected := mockedResponse{
    body: `{"thing": "thingamabob", "answer": 42}`,
    cType: "application/json",
    status: 200,
  }

  res, _ := SuccessResponse(wadus, mockReq)
  if res.StatusCode != expected.status {
    t.Errorf("expected status=%d, got %d", expected.status, res.StatusCode)
  }
  if v, ok := res.Headers["Content-Type"]; !ok || v != expected.cType {
    t.Errorf("expected content-type=%s, got '%s'", expected.cType, v)
  }

  ja.Assertf(string(res.Body), expected.body)
}

