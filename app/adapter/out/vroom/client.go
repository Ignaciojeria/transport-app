package vroom

import (
	"net/http"
	"time"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/hashicorp/go-retryablehttp"
)

func init() {
	ioc.Registry(NewVroomFastClient)
	ioc.Registry(NewVroomDefaultClient)
	ioc.Registry(NewVroomHeavyClient)
}

func NewVroomFastClient() *retryablehttp.Client {
	return buildClient(2 * time.Minute)
}

func NewVroomDefaultClient() *retryablehttp.Client {
	return buildClient(5 * time.Minute)
}

func NewVroomHeavyClient() *retryablehttp.Client {
	return buildClient(10 * time.Minute)
}

func buildClient(timeout time.Duration) *retryablehttp.Client {
	c := retryablehttp.NewClient()
	c.RetryMax = 3
	c.RetryWaitMin = 1 * time.Second
	c.RetryWaitMax = 5 * time.Second
	c.HTTPClient.Timeout = timeout
	c.HTTPClient.Transport = &http.Transport{
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	c.Logger = nil
	return c
}
