package main

import (
  "context"
)


type IUser interface {
  Username() string
}

func GetBirthday(ctx context.Context, user IUser) (*Birthday, error) {
  bday := newBirthday(user.Username())
  err := GetFromStore(&bday)
  return &bday, err
}

