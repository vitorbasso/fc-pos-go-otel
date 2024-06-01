package gateway

import "context"

type TemperatureResponse struct {
	Celsius    float64
	Fahrenheit float64
}

type TemperatureGateway interface {
	FetchTemperatureByCity(ctx context.Context, city string) (*TemperatureResponse, error)
}
