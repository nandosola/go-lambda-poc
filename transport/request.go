package transport


import (
  "context"
  "fmt"

  "github.com/go-playground/validator/v10"
)


const userNameParam = "username"

var validate *validator.Validate

func init() {
  validate = validator.New(validator.WithRequiredStructEnabled())
}

type GetBirthdayRequest struct {
  Name  string `validate:"required,alpha,lowercase,min=2,max=12"`
}

func NewReadRequest(ctx context.Context, req *Request) (*GetBirthdayRequest, error) {
  name, ok := req.PathParameters[userNameParam]
  if !ok {
    return nil, fmt.Errorf("PathParamNotFound: '%s'", userNameParam)
  }

  br := GetBirthdayRequest{
    Name: name,
  }

  if err := validate.Struct(&br); err != nil {
    return nil, err
  }

  return &br, nil
}


func (gbr *GetBirthdayRequest) Username() string {
  return gbr.Name
}

