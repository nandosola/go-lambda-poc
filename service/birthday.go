package service

import (
  "crypto/sha256"
  "encoding/json"
  "fmt"
  "time"

  "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
  "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)


const (
  defaultGreetingTmpl = "Hello, %s! Your birthday is in %d day(s)"
  bdayGreetingTmpl    = "Hello, %s! Happy birthday!"
)

type Birthday struct {
  Id        string    `dynamodbav:"Id"   json:"-"`
  Dob       time.Time `dynamodbav:"Dob"  json:"-"`
  name      string
}

func newBirthday(userName string) Birthday {
  hash := sha256.Sum256([]byte(userName))
  return Birthday{
    Id: fmt.Sprintf("%x", hash[:]),
    name: userName,
  }
}

func (b Birthday) daysRemaining() uint {
  today := time.Now().UTC()

  ty, tm, td := today.Date()
  _, bm, bd := b.Dob.Date()

  if tm == bm && td == bd {
    return 0
  }

  var nextBday time.Time
  if tm > bm || (tm == bm && td > bd) {
    nextBday = time.Date(ty+1, time.Month(bm), bd, 0, 0, 0, 0, time.UTC)
  } else {
    nextBday = time.Date(ty, time.Month(bm), bd, 0, 0, 0, 0, time.UTC)
  }
  days := nextBday.Sub(today).Hours() / 24

  return uint(days)
}

func (b Birthday) greet() string {
  days := b.daysRemaining()
  if 0 == days {
    return fmt.Sprintf(bdayGreetingTmpl, b.name)
  }

  return fmt.Sprintf(defaultGreetingTmpl, b.name, days)
}

// Define serializers/deserializers/views in the same struct. DTOs are not idiomatic.

func (b Birthday) GetKey() map[string]types.AttributeValue {
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

