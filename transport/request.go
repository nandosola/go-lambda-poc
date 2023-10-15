package transport


import (
  "context"
  "fmt"
  "strings"

  "github.com/go-playground/validator/v10"
)


const (
  userNameParam       = "username"
  readReqRestrictions = "alphanumeric, min=2, max=12"
)

var (
  validate *validator.Validate
)

func init() {
  validate = validator.New(validator.WithRequiredStructEnabled())
}

type GetBirthdayRequest struct {
  Name  string `validate:"required,alpha,min=2,max=12"`  // tags must be literals
}

func NewReadRequest(ctx context.Context, req *Request) (*GetBirthdayRequest, error) {
  name, ok := req.PathParameters[userNameParam]
  if !ok {
    return nil, fmt.Errorf("PathParamNotFound: '%s'", userNameParam)
  }

  br := GetBirthdayRequest{
    Name: strings.ToLower(name),
  }

  if err := validate.Struct(&br); err != nil {
    return nil, fmt.Errorf("ValidationError.Read: %w", err)
  }

  return &br, nil
}

func (gbr *GetBirthdayRequest) Username() string {
  return gbr.Name
}

