package service

import (
  "context"
  "log"
  "sync"
)


var (
  readRepo IBirthdayReader
  onceRepo sync.Once
)


type IUser interface {
  Username() string
}

type IBirthdayStore interface {
  GetFromStore(*Birthday) error
}

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
  bday := newBirthday(user.Username())
  err := brr.store.GetFromStore(&bday)
  return &bday, err
}

