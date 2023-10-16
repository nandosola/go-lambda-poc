package service

import (
  "crypto/sha256"
  "encoding/json"
  "errors"
  "fmt"
  "time"

  "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
  "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)


const (
  defaultGreetingTmpl = "Hello, %s! Your birthday is in %d day(s)"
  bdayGreetingTmpl    = "Hello, %s! Happy birthday!"
)

// package-level clock impl
type timeNowFunT func() time.Time

var (
  NowFun timeNowFunT
  ErrInvalidBirthday = errors.New("birth date must be before today")
)

func init() {
  ResetClock()
}

func ResetClock() {
  NowFun = func() time.Time {
    return time.Now()
  }
}

type Birthday struct {
  Id        string    `dynamodbav:"Id"   json:"-"`
  Dob       time.Time `dynamodbav:"Dob"  json:"-"`
  name      string
}

func NewBirthday(userName string) *Birthday {
  hash := sha256.Sum256([]byte(userName))
  return &Birthday{
    Id: fmt.Sprintf("%x", hash[:]),
    name: userName,
  }
}

func (b *Birthday) WithDateOfBirth(dob time.Time) (*Birthday, error) {
  now := NowFun().UTC()

  ty, tm, td := now.Date()
  by, bm, bd := dob.Date()


  if (ty == by && tm == bm && td == bd) || !dob.Before(now) {
   return nil, fmt.Errorf("Birthday: %w", ErrInvalidBirthday)
  }

  b.Dob = dob
  return b, nil
}

func (b *Birthday) daysRemaining() uint {
  now := NowFun().UTC()
  today := now.YearDay()
  bday := b.Dob.YearDay()

  if today <= bday {
    return uint(bday-today)
  }

  // bday is next year
  ty, _, _ := now.Date()
  _, bm, bd := b.Dob.Date()
  end := time.Date(ty, 12, 31, 0, 0, 0, 0, time.UTC).YearDay()  // == 366 on leap years
  next := time.Date(ty+1, bm, bd, 0, 0, 0, 0, time.UTC).YearDay()

  return uint(end-today+next)
}

func (b *Birthday) greet() string {
  days := b.daysRemaining()
  if 0 == days {
    return fmt.Sprintf(bdayGreetingTmpl, b.name)
  }

  return fmt.Sprintf(defaultGreetingTmpl, b.name, days)
}

// Define serializers/deserializers/views in the same struct. DTOs are not idiomatic.

func (b *Birthday) GetKey() map[string]types.AttributeValue {
  id, err := attributevalue.Marshal(b.Id)
  if err != nil {
    panic(err)
  }
  return map[string]types.AttributeValue{"Id": id}
}

func (b Birthday) MarshalJSON() ([]byte, error) {
  return json.Marshal(&struct {
    Message  string  `json:"message"`
  }{
    Message: b.greet(),
  })
}

