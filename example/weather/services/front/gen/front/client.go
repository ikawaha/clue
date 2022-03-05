// Code generated by goa v3.6.0, DO NOT EDIT.
//
// front client
//
// Command:
// $ goa gen goa.design/clue/example/weather/services/front/design -o
// services/front

package front

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Client is the "front" service client.
type Client struct {
	ForecastEndpoint goa.Endpoint
}

// NewClient initializes a "front" service client given the endpoints.
func NewClient(forecast goa.Endpoint) *Client {
	return &Client{
		ForecastEndpoint: forecast,
	}
}

// Forecast calls the "forecast" endpoint of the "front" service.
// Forecast may return the following errors:
//	- "not_usa" (type *goa.ServiceError): IP address is not in the US
//	- error: internal error
func (c *Client) Forecast(ctx context.Context, p string) (res *Forecast2, err error) {
	var ires interface{}
	ires, err = c.ForecastEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*Forecast2), nil
}
