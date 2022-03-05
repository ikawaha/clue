// Code generated by goa v3.6.0, DO NOT EDIT.
//
// locator gRPC client types
//
// Command:
// $ goa gen goa.design/clue/example/weather/services/locator/design -o
// services/locator

package client

import (
	locatorpb "goa.design/clue/example/weather/services/locator/gen/grpc/locator/pb"
	locator "goa.design/clue/example/weather/services/locator/gen/locator"
)

// NewGetLocationRequest builds the gRPC request type from the payload of the
// "get_location" endpoint of the "locator" service.
func NewGetLocationRequest(payload string) *locatorpb.GetLocationRequest {
	message := &locatorpb.GetLocationRequest{}
	message.Field = payload
	return message
}

// NewGetLocationResult builds the result type of the "get_location" endpoint
// of the "locator" service from the gRPC response type.
func NewGetLocationResult(message *locatorpb.GetLocationResponse) *locator.WorldLocation {
	result := &locator.WorldLocation{
		Lat:     message.Lat,
		Long:    message.Long,
		City:    message.City,
		Region:  message.Region,
		Country: message.Country,
	}
	return result
}
