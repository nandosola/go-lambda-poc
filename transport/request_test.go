package transport

import (
  "context"
  "encoding/json"
  "errors"
  "fmt"
  "strings"
  "testing"
  "time"

  "service"
)


func TestGetBirthdayRequestEmptyPathParams(t *testing.T) {
  mockReq := &Request{ PathParameters: make(map[string]string)}
  _, err := NewReadRequest(context.TODO(), mockReq)
  if !errors.Is(err, ErrPathParamNotFound) {
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

func TestAddOrUpdateBirthdayRequest(t *testing.T) {
  defer service.ResetClock()
  service.NowFun = func() time.Time {
      return time.Date(1823, 8, 3, 13, 37, 42, 0, time.UTC)
  }

  cases := []struct {
    name      string
    user      string
    birthDate string
    ok        bool
  }{
    {
      name: "SetInThePast",
      user: "ludwigvan",
      birthDate: "1770-12-17",
      ok: true,
    },
    {
      name: "SetInTheFuture",
      user: "gusmahler",
      birthDate: "1860-07-07",
      ok: false,
    },
    {
      name: "SetToday",
      user: "fbarbieri",
      birthDate: "1823-08-03",
      ok: false,
    },
    {
      name: "Empty",
      user: "",
      birthDate: "1808-05-02",
      ok: false,
    },
    {
      name: "MinLenFail",
      user: "x",
      birthDate: "1808-05-02",
      ok: false,
    },
    {
      name: "MaxLenFail1",
      user: "foobarbazbatq",
      birthDate: "1808-05-02",
      ok: false,
    },
    {
      name: "MaxLenFail2",
      user: "foobarbazbatqux",
      birthDate: "1808-05-02",
      ok: false,
    },
    {
      name: "InvalidChars1",
      user: "kk123456",
      birthDate: "1808-05-02",
      ok: false,
    },
    {
      name: "InvalidChars2",
      user: "alice.json",
      birthDate: "1808-05-02",
      ok: false,
    },
    {
      name: "InvalidChars3",
      user: "Robert%60%29%3B%20DROP%20TABLE%20Students%3B--",
      birthDate: "1808-05-02",
      ok: false,
    },
  }

  for _, tc := range cases {
    t.Run(tc.name, func(t *testing.T) {
      mockReq := &Request{
        PathParameters: map[string]string{ "username": tc.user },
        Body: fmt.Sprintf(`{"dateOfBirth":"%s"}`, tc.birthDate),
      }

      _, err := NewAddOrUpdateRequest(context.TODO(), mockReq)
      if err != nil && tc.ok {
        t.Errorf("%s was expected to pass. error: %s", tc.user, err.Error())
      }
    })
  }
}

func TestAddOrUpdateBirthdayJSONUnmarshal(t *testing.T) {
  cases := []struct {
    name  string
    input string
    ok    bool
  }{
    {
      name: "BadJSON",
      input: "not_json" ,
      ok: false,
    },
    {
      name: "EmptyObject",
      input: "{}" ,
      ok: false,
    },
    {
      name: "MissingKey",
      input: `{"fruit":"pineapple"}`,
      ok: false,
    },
    {
      name: "WrongKeyType",
      input: `{"dateOfBirth":42}`,
      ok: false,
    },
    {
      name: "WrongKeyFmt",
      input: `{"dateOfBirth":"this-is-not-a-date"}`,
      ok: false,
    },
    {
      name: "BadDate",
      input: `{"dateOfBirth":"1492-1-2"}`,
      ok: false,
    },
    {
      name: "GoodDatePast",
      input: `{"dateOfBirth":"1492-01-02"}`,
      ok: true,
    },
    {
      name: "GoodDateFuture",  // NB: this date will get rejected by service.Birthday invariants
      input: `{"dateOfBirth":"2718-03-14"}`,
      ok: true,
    },
  }

  for _, tc := range cases {
    t.Run(tc.name, func(t *testing.T) {
      var req AddOrUpdateBirthdayRequest
      if err := json.Unmarshal([]byte(tc.input), &req); err != nil && tc.ok {
        t.Errorf("unexpected error: %s", err)
      }
    })
  }
}

