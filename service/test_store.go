package service

import (
  "encoding/json"
  "io"
  "sync"
  "time"
)

var (
  TestRWStore   *testStore
  onceTestStore sync.Once
)


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
  onceTestStore.Do(func(){
    TestRWStore = &testStore{birthdays: make(map[string]time.Time)}
  })
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

func (ts *testStore) GetFromStore(bday *Birthday) (bool, error) {
  v, ok := TestRWStore.birthdays[bday.Id]
  if ok {
    bday.Dob = v
    return true, nil
  }
  return false, nil
}

func (ds *testStore) AddToStore(bday *Birthday) error {
 TestRWStore.birthdays[bday.Id] = bday.Dob
 return nil
}

func (ts *testStore) Clean(){
  m := TestRWStore.birthdays
  for k := range m {
    delete(m, k)
  }
}

