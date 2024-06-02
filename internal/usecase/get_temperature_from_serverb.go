package usecase

import (
	"context"
	"fmt"
	"otel/internal/constants"
	"otel/internal/gateway"

	"go.opentelemetry.io/otel/trace"
)

type GetTemperatureFromServerBResponse struct {
	City       string
	Celsius    string
	Fahrenheit string
	Kelvin     string
}

func NewGetTemperatureFromServerBUseCase(temperatureFromCepGateway gateway.TemperatureByCepGateway) *GetTemperatureFromServerBUseCase {
	return &GetTemperatureFromServerBUseCase{
		temperatureFromCepGateway: temperatureFromCepGateway,
	}
}

type GetTemperatureFromServerBUseCase struct {
	temperatureFromCepGateway gateway.TemperatureByCepGateway
}

func (g *GetTemperatureFromServerBUseCase) Execute(ctx context.Context, cep string) (*GetTemperatureFromServerBResponse, error) {
	tracer := ctx.Value(constants.CtxTracerKey).(trace.Tracer)
	ctx, span := tracer.Start(ctx, "get_temperature_from_serverb_use_case.Execute")
	defer span.End()

	temperature, err := g.temperatureFromCepGateway.FetchTemperatureByCep(ctx, cep)
	if err != nil {
		return nil, fmt.Errorf("get_temperature_from_serverb_use_case.Execute: %w", err)
	}
	return &GetTemperatureFromServerBResponse{
		City:       temperature.City,
		Celsius:    temperature.Celsius,
		Fahrenheit: temperature.Fahrenheit,
		Kelvin:     temperature.Kelvin,
	}, nil
}
