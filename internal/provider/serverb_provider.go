package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"otel/internal/gateway"
	"time"
)

var (
	ErrUnexpectedResponse = fmt.Errorf("unexpected response")
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
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	url := fmt.Sprintf("http://%s:%s/temperatures/%s", provider.serverBHost, provider.serverBPort, cep)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("serverb_provider.GetTemperatureFromCep: %w", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("serverb_provider.GetTemperatureFromCep: %w", err)
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
