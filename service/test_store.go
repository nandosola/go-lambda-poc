package service

import (
  "encoding/json"
  "io"
  "time"
)

var TestRWStore *testStore

type testStore struct {
  birthdays map[string]time.Time
}

type Fixture struct {
  Birthdays []struct {
    PutRequest struct {
      Item struct {
        Id struct {
          S string `json:"S"`
        } `json:"Id"`
        Dob struct {
          S time.Time `json:"S"`
        } `json:"Dob"`
      } `json:"Item"`
    } `json:"PutRequest"`
  } `json:"Birthdays"`
}

func testConnect() *testStore {
  TestRWStore = &testStore{birthdays: make(map[string]time.Time)}
  return TestRWStore
}

func (ts *testStore) Load(jsonFixture io.Reader) error {
  var fxt Fixture
  decoder := json.NewDecoder(jsonFixture)
  if err := decoder.Decode(&fxt); err != nil {
    return err
  }

  for _, req := range fxt.Birthdays {
    item := req.PutRequest.Item
    TestRWStore.birthdays[item.Id.S] = item.Dob.S
  }

  return nil
}


func (ts *testStore) GetFromStore(bday *Birthday) error {
  bday.Dob = TestRWStore.birthdays[bday.Id]
  return nil
}

func (ts *testStore) Clean(){
}

