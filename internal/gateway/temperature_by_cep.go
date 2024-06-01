package gateway

import "context"

type TemperatureByCepResponse struct {
	City       string
	Celsius    string
	Fahrenheit string
	Kelvin     string
}

type TemperatureByCepGateway interface {
	FetchTemperatureByCep(ctx context.Context, cep string) (*TemperatureByCepResponse, error)
}
