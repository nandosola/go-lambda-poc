package transport

import (
  "context"
  "strings"
  "testing"
)


func TestGetBirthdayRequestEmptyPathParams(t *testing.T) {
  mockReq := &Request{ PathParameters: make(map[string]string)}
  _, err := NewReadRequest(context.TODO(), mockReq)
  if !strings.HasPrefix(err.Error(), "PathParamNotFound") {
    t.Error("expected PathParamNotFound")
  }
}

func TestGetBirthdayRequest(t *testing.T) {
  cases := []struct {
    name string
    user string
    ok   bool
  }{
    {
      name: "Pass1",
      user: "alice",
      ok: true,
    },
    {
      name: "Pass2",
      user: "al",
      ok: true,
    },
    {
      name: "Pass3",
      user: "foobarbazbat",
      ok: true,
    },
    {
      name: "Empty",
      user: "",
      ok: false,
    },
    {
      name: "MinLenFail",
      user: "x",
      ok: false,
    },
    {
      name: "MaxLenFail1",
      user: "foobarbazbatq",
      ok: false,
    },
    {
      name: "MaxLenFail2",
      user: "foobarbazbatqux",
      ok: false,
    },
    {
      name: "InvalidChars1",
      user: "kk123456",
      ok: false,
    },
    {
      name: "InvalidChars2",
      user: "alice.json",
      ok: false,
    },
    {
      name: "InvalidChars3",
      user: "Robert%60%29%3B%20DROP%20TABLE%20Students%3B--",
      ok: false,
    },
  }

  for _, tc := range cases {
    t.Run(tc.name, func(t *testing.T) {
      mockReq := &Request{ PathParameters: map[string]string{ "username": tc.user }}

      _, err := NewReadRequest(context.TODO(), mockReq)
      if err != nil && tc.ok {
        t.Errorf("%s was expected to pass. error: %s", tc.user, err.Error())
      }
      if !tc.ok && !strings.HasPrefix(err.Error(), "ValidationError") {
        t.Errorf("%s: expected ValidationError", tc.user)
      }
    })
  }
}

