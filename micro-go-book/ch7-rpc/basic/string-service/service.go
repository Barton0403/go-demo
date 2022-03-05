package service

type StringRequest struct {
  A string
  B string
}

type Service interface {
  Concat(req StringRequest, ret *string) error

  Diff(req StringRequest, ret *string) error
}

type StringService struct {
}

func (s StringService) Concat(req StringRequest, ret *string) error {
  *ret = "1"
  return nil
}

func (s StringService) Diff(req StringRequest, ret *string) error {
  *ret = "2"
  return nil
}
