package main


import (
  "context"

  "github.com/go-playground/validator/v10"
)


var validate *validator.Validate

func init() {
  validate = validator.New(validator.WithRequiredStructEnabled())
}

type GetBirthdayRequest struct {
  Name  string `validate:"required,alpha,lowercase,min=2,max=12"`
}

func newReadRequest(ctx context.Context, name string) (*GetBirthdayRequest, error) {
  req := GetBirthdayRequest{
    Name: name,
  }

  if err := validate.Struct(&req); err != nil {
    return nil, err
  }

  return &req, nil
}


func (gbr *GetBirthdayRequest) Username() string {
  return gbr.Name
}

