package service

import (
  "bytes"
  "context"
  "fmt"
  "io"
  "os"
  "testing"
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
  LoadFixture(t, "../testdata/fixtures.json")
  alice :=  TestUser{name:"alice"}
  res, err := Reader().GetBirthday(context.TODO(), alice)
  if err != nil {
    t.Fatal(err)
  }

  fmt.Printf("%+v\n", res)

}
