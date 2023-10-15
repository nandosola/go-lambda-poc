package service

import (
  "context"
  "errors"
  "fmt"
  "log"
  "sync"
)


var (
  readRepo IBirthdayReader
  onceRepo sync.Once

  ErrNotFound = errors.New("username not found")
)

// Decouple request objects from transport layer through interfaces
type IUser interface {
  Username() string
}

type IBirthdayStore interface {
  GetFromStore(*Birthday) (bool, error)
}

// Consume from GET
type IBirthdayReader interface {
  GetBirthday(context.Context, IUser) (*Birthday, error)
}

type bdayReadRepository struct {
  store IBirthdayStore
}

func Reader() IBirthdayReader {
  if readRepo == nil {
    log.Fatal("BirthdayReader is not initialized")
  }
  return readRepo
}

func InitializeTestRepo() error {
  onceRepo.Do(func(){
    readRepo = bdayReadRepository{
      store: testConnect(),
    }
  })

  return nil
}

func InitializeRepo() error {
  st, err := ddbConnect()
  if err != nil {
    return err
  }

  onceRepo.Do(func(){
    readRepo = bdayReadRepository{
      store: st,
    }
  })

  return nil
}

func (brr bdayReadRepository) GetBirthday(ctx context.Context, user IUser) (*Birthday, error) {
  u := user.Username()
  bday := newBirthday(u)
  ok, err := brr.store.GetFromStore(&bday)
  if !ok {
    return nil, fmt.Errorf("GetBirthday: %s, %w", u, ErrNotFound)
  }
  return &bday, err
}

