package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"otel/internal/constants"
	"otel/internal/gateway"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrUnexpectedResponse = fmt.Errorf("unexpected response")
	ErrNotFound           = fmt.Errorf("cep not found")
)

type serverBProvider struct {
	serverBHost string
	serverBPort string
}

type serverBResponse struct {
	City       string `json:"city"`
	Celsius    string `json:"temp_C"`
	Fahrenheit string `json:"temp_F"`
	Kelvin     string `json:"temp_K"`
}

func NewServerBProvider(serverBHost, serverBPort string) *serverBProvider {
	return &serverBProvider{
		serverBHost: serverBHost,
		serverBPort: serverBPort,
	}
}

func (provider *serverBProvider) FetchTemperatureByCep(ctx context.Context, cep string) (*gateway.TemperatureByCepResponse, error) {
	tracer := ctx.Value(constants.CtxTracerKey).(trace.Tracer)
	ctx, span := tracer.Start(ctx, "serverb_provider.FetchTemperatureByCep")
	defer span.End()

	url := fmt.Sprintf("http://%s:%s/temperatures/%s", provider.serverBHost, provider.serverBPort, cep)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("serverb_provider.GetTemperatureFromCep: %w", err)
	}
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("serverb_provider.GetTemperatureFromCep: %w", err)
	}
	if res.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("serverb_provider.GetTemperatureFromCep: %w", ErrNotFound)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("serverb_provider.GetTemperatureFromCep: %w; status code %d", ErrUnexpectedResponse, res.StatusCode)
	}
	defer res.Body.Close()
	var response serverBResponse
	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("serverb_provider.GetTemperatureFromCep: %w", err)
	}
	return &gateway.TemperatureByCepResponse{
		City:       response.City,
		Celsius:    response.Celsius,
		Fahrenheit: response.Fahrenheit,
		Kelvin:     response.Kelvin,
	}, nil
}
