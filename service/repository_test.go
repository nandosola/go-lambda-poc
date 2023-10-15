package service

import (
  "bytes"
  "context"
  "errors"
  "fmt"
  "io"
  "os"
  "testing"
  "time"
)

func TestMain(m *testing.M) {
  InitializeTestRepo()
  exitVal := m.Run()
  os.Exit(exitVal)
}

func LoadFixture(t *testing.T, file string) {
  data, err := os.Open(file)
  if err != nil {
    t.Fatal(err)
  }

  loadedData, err := io.ReadAll(data)
  if err != nil {
    return
  }

  bufUsers := new(bytes.Buffer)
  bufUsers.Write(loadedData)
  if err = TestRWStore.Load(bufUsers); err != nil {
    t.Fatal(err)
  }
  if err = data.Close(); err != nil {
    t.Fatal(err)
  }
  t.Cleanup(func() {
    fmt.Println("cleanup")
    TestRWStore.Clean()
  })
}

type TestUser struct {
  name string
}

func (tu TestUser) Username() string {
  return tu.name
}

func TestGetBirthday(t *testing.T) {
  defer resetClock()
  nowFun = func() time.Time {
      return time.Date(2022, 6, 18, 13, 37, 42, 0, time.UTC)
  }

  // shared fixture
  LoadFixture(t, "../testdata/fixtures.json")

  cases := []struct {
    name   string
    user   string
    found  bool
  }{
    {
      name: "BdayToday",
      user: "alice",
      found: true,
    },
    {
      name: "BdayThisYear",
      user: "bob",
      found: true,
    },
    {
      name: "BdayNextYear",
      user: "charly",
      found: true,
    },
    {
      name: "NotFound",
      user: "bogus",
      found: false,
    },
  }

  for _, tc := range cases {
    t.Run(tc.name, func(t *testing.T) {
      user :=  TestUser{name: tc.user}
      _, err := Reader().GetBirthday(context.TODO(), user)
      if err != nil && tc.found {
        t.Fatal(err)
      }
      if !tc.found && !errors.Is(err, ErrNotFound) {
        t.Errorf("%s: expected ErrNotFound", tc.user)
      }
    })
  }
}

