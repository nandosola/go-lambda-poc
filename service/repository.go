package service

import (
  "context"
  "errors"
  "fmt"
  "log"
  "sync"
)


var (
  readRepo  IBirthdayReader
  writeRepo IBirthdayWriter
  onceRepo  sync.Once

  ErrNotFound = errors.New("username not found")
)

// Decouple request objects from transport layer through interfaces
type IUser interface {
  Username() string
}

type IBirthdayStore interface {
  GetFromStore(*Birthday) (bool, error)
  AddToStore(*Birthday) error
}

// Consume from GET
type IBirthdayReader interface {
  GetBirthday(context.Context, IUser) (*Birthday, error)
}

// Consume from PUT
type IBirthdayWriter interface {
  AddBirthday(context.Context, *Birthday) error
}

type bdayReadRepository struct {
  store IBirthdayStore
}

type bdayWriteRepository struct {
  store IBirthdayStore
}

func Reader() IBirthdayReader {
  if readRepo == nil {
    log.Fatal("BirthdayReader is not initialized")
  }
  return readRepo
}

func Writer() IBirthdayWriter {
  if writeRepo == nil {
    log.Fatal("BirthdayWriter is not initialized")
  }
  return writeRepo
}

func InitializeTestRepo() error {
  onceRepo.Do(func(){
    readRepo = bdayReadRepository{
      store: testConnect(),
    }
    writeRepo = bdayWriteRepository{
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
    writeRepo = bdayWriteRepository{
      store: st,
    }
  })

  return nil
}

//
// About context handling:
// lambda ctx is not of much use in our case. It should be handled by the repo though.
//   - extract auth metadata (eg Cognito
//   - handle long-running tasks and timeout deadlines
//   see: https://docs.aws.amazon.com/lambda/latest/dg/golang-context.html

func (brr bdayReadRepository) GetBirthday(ctx context.Context, user IUser) (*Birthday, error) {
  u := user.Username()
  bday := NewBirthday(u)
  ok, err := brr.store.GetFromStore(bday)
  if !ok {
    return nil, fmt.Errorf("GetBirthday: %s, %w", u, ErrNotFound)
  }
  return bday, err
}

// no need to test this
func (bwr bdayWriteRepository) AddBirthday(ctx context.Context, bday *Birthday) error {
  return bwr.store.AddToStore(bday)
}

