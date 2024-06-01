package di

import (
	"otel/configs"
	"otel/internal/infra/webserver/handler"
	"otel/internal/provider"
	"otel/internal/usecase"
)

func NewTemperatureHandler(viacepUrl, weatherApiUrl, weatherApiKey string) *handler.TemperatureHandler {
	config := configs.GetConfig()
	locationGateway := provider.NewViaCepLocationProvider(viacepUrl)
	temperatureGateway := provider.NewWeatherAPIProvider(weatherApiKey, weatherApiUrl)
	temperatureByCepGateway := provider.NewServerBProvider(config.ServerBHost, config.ServerBPort)
	usecaseB := usecase.NewGetTemperatureFromCepUseCase(locationGateway, temperatureGateway)
	usecaseA := usecase.NewGetTemperatureFromServerBUseCase(temperatureByCepGateway)
	return handler.NewTemperatureHandler(usecaseB, usecaseA)
}
