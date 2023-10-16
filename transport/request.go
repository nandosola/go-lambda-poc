package transport


import (
  "context"
  "encoding/json"
  "errors"
  "fmt"
  "regexp"
  "strings"
  "time"

  "service"

  "github.com/go-playground/validator/v10"
)


const (
  userNameParam         = "username"
  readReqRestrictions   = "alphanumeric, min=2, max=12"
  updateReqRestrictions = "yyyy-mm-dd"

  yyyymmdd            = "2006-01-02"
)

var (
  validate *validator.Validate

  ErrPathParamNotFound = errors.New("missing path param")
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
    return nil, fmt.Errorf("NewReadRequest: %s, %w", userNameParam, ErrPathParamNotFound)
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

type AddOrUpdateBirthdayRequest struct {
  Dob       string    `json:"dateOfBirth" validate:"required,datetime"`
  dobAsDate time.Time
}

func NewAddOrUpdateRequest(ctx context.Context, req *Request) (*service.Birthday, error) {
  gbr, err := NewReadRequest(ctx, req)  // reuse path param validation
  if err != nil {
    return nil, err
  }

  var update AddOrUpdateBirthdayRequest
  if err := json.Unmarshal([]byte(req.Body), &update); err != nil {
    return nil, fmt.Errorf("JSONUnmarshal: %w", err)
  }

  return service.NewBirthday(gbr.Name).WithDateOfBirth(update.dobAsDate)
}

func (aur *AddOrUpdateBirthdayRequest) UnmarshalJSON(data []byte) error {
  type Alias AddOrUpdateBirthdayRequest
  a := (*Alias)(aur)

  if err := json.Unmarshal(data, &a); err != nil {
    return err
  }

  if aur.Dob == "" {
   return errors.New("ValidationError.Update: 'dateOfBirth' is required",)
  }

  bday, err := asDate(aur.Dob)
  if err != nil {
    return err
  }

  aur.dobAsDate = bday

  return nil
}

func asDate(str string) (time.Time, error){
  re := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
  if !re.MatchString(str) {
   return time.Time{}, errors.New("ValidationError.Update: bad date format")
  }

  parsed, err := time.Parse(yyyymmdd, str)
  if err != nil {
   return time.Time{}, fmt.Errorf("ValidationError.Update: %w", err)
  }

  return parsed, nil
}

