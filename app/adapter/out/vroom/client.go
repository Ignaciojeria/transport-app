package vroom

import (
	"time"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-resty/resty/v2"
)

func init() {
	ioc.Registry(NewVroomRestyFastClient)
	ioc.Registry(NewVroomRestyDefaultClient)
	ioc.Registry(NewVroomRestyHeavyClient)
}

func NewVroomRestyFastClient() *resty.Client {
	return buildRestyClient(2 * time.Minute)
}

func NewVroomRestyDefaultClient() *resty.Client {
	return buildRestyClient(5 * time.Minute)
}

func NewVroomRestyHeavyClient() *resty.Client {
	return buildRestyClient(10 * time.Minute)
}

func buildRestyClient(timeout time.Duration) *resty.Client {
	client := resty.New()
	client.SetTimeout(timeout)

	// Puedes personalizar más abajo si quieres:
	client.
		SetRetryCount(0).
		SetHeader("Content-Type", "application/json")

	return client
}
