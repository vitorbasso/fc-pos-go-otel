package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"otel/internal/constants"
	"otel/internal/gateway"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrWeatherAPIStatusNotOk = errors.New("weather api status not ok")
)

type weatherApiResponse struct {
	Current struct {
		LastUpdated string  `json:"last_updated"`
		TempC       float64 `json:"temp_c"`
		TempF       float64 `json:"temp_f"`
	} `json:"current"`
}

type WeatherAPIProvider struct {
	WeatherAPIKey string
	WeatherAPIUrl string
}

func NewWeatherAPIProvider(weatherAPIKey, weatherAPIUrl string) *WeatherAPIProvider {
	return &WeatherAPIProvider{
		WeatherAPIKey: weatherAPIKey,
		WeatherAPIUrl: weatherAPIUrl,
	}
}

func (w *WeatherAPIProvider) FetchTemperatureByCity(ctx context.Context, city string) (*gateway.TemperatureResponse, error) {
	tracer := ctx.Value(constants.CtxTracerKey).(trace.Tracer)
	ctx, span := tracer.Start(ctx, "weather_api_provider.FetchTemperatureByCity")
	defer span.End()

	url := fmt.Sprintf("%s?key=%s&q=%s&aqi=no", w.WeatherAPIUrl, w.WeatherAPIKey, url.QueryEscape(city))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("weather_api_provider.FetchTemperatureByCep: %w", err)
	}
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("weather_api_provider.FetchTemperatureByCep: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather_api_provider.FetchTemperatureByCep: %w; code: %s", ErrWeatherAPIStatusNotOk, res.Status)
	}
	var response weatherApiResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("weather_api_provider.FetchTemperatureByCep: %w", err)
	}
	return &gateway.TemperatureResponse{
		Celsius:    response.Current.TempC,
		Fahrenheit: response.Current.TempF,
	}, nil
}
