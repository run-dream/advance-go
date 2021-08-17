package service

import (
	"context"
	"go-advance/src/concurrency/api/v1"
)

type HelloServiceImpl struct {
}

func (s *HelloServiceImpl) Hello(c context.Context, str *api.String) (*api.String, error) {
	return &api.String{Value: "hello" + str.Value}, nil
}
