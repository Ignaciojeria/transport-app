package usecase

import (
	"context"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type Signup func(ctx context.Context, input interface{}) (interface{}, error)

func init() {
	ioc.Registry(NewSignup)
}

func NewSignup() Signup {
	return func(ctx context.Context, input interface{}) (interface{}, error) {
		return input, nil
	}
}
