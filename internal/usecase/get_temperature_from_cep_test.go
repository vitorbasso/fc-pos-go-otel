package usecase

import (
	"context"
	"fmt"
	"otel/internal/gateway"
	"otel/internal/provider"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockLocationGateway struct {
	mock.Mock
}

func (m *mockLocationGateway) FetchLocationByCep(ctx context.Context, cep string) (*gateway.LocationResponse, error) {
	args := m.Called(ctx, cep)
	return args.Get(0).(*gateway.LocationResponse), args.Error(1)
}

type mockTemperatureGateway struct {
	mock.Mock
}

func (m *mockTemperatureGateway) FetchTemperatureByCity(ctx context.Context, city string) (*gateway.TemperatureResponse, error) {
	args := m.Called(ctx, city)
	return args.Get(0).(*gateway.TemperatureResponse), args.Error(1)
}

func TestGetTemperatureFromCepResponse(t *testing.T) {
	locationMock := new(mockLocationGateway)
	temperatureMock := new(mockTemperatureGateway)

	cep := "12345678"
	locationMock.On("FetchLocationByCep", mock.Anything, cep).Return(&gateway.LocationResponse{
		City: "São Paulo",
	}, nil)
	temperatureMock.On("FetchTemperatureByCity", mock.Anything, "São Paulo").Return(&gateway.TemperatureResponse{
		Celsius: 25.0,
	}, nil)

	usecase := NewGetTemperatureFromCepUseCase(locationMock, temperatureMock)
	usecaseResponse, err := usecase.Execute(context.Background(), cep)

	if assert.Nil(t, err) {
		assert.Equal(t, "São Paulo", usecaseResponse.City)
		assert.Equal(t, "25.00", usecaseResponse.Celsius)
		assert.Equal(t, "77.00", usecaseResponse.Fahrenheit)
		assert.Equal(t, "298.00", usecaseResponse.Kelvin)
	}

	locationMock.AssertExpectations(t)
	temperatureMock.AssertExpectations(t)
}

func TestGetTemperatureFromCepResponse_CepNotFound(t *testing.T) {
	locationMock := new(mockLocationGateway)
	temperatureMock := new(mockTemperatureGateway)

	cep := "12345678"
	locationMock.On("FetchLocationByCep", mock.Anything, cep).Return((*gateway.LocationResponse)(nil), gateway.ErrCepNotFound)

	usecase := NewGetTemperatureFromCepUseCase(locationMock, temperatureMock)
	usecaseResponse, err := usecase.Execute(context.Background(), cep)
	assert.ErrorIs(t, err, ErrNotFound)
	assert.Nil(t, usecaseResponse)

	locationMock.AssertExpectations(t)
	temperatureMock.AssertNotCalled(t, "FetchTemperatureByCity")
}

func TestGetTemperatureFromCepResponse_TemperatureGatewayError(t *testing.T) {
	locationMock := new(mockLocationGateway)
	temperatureMock := new(mockTemperatureGateway)

	cep := "12345678"
	locationMock.On("FetchLocationByCep", mock.Anything, cep).Return(&gateway.LocationResponse{
		City: "São Paulo",
	}, nil)
	temperatureMock.On("FetchTemperatureByCity", mock.Anything, "São Paulo").Return((*gateway.TemperatureResponse)(nil), provider.ErrWeatherAPIStatusNotOk)

	usecase := NewGetTemperatureFromCepUseCase(locationMock, temperatureMock)
	usecaseResponse, err := usecase.Execute(context.Background(), cep)
	assert.ErrorIs(t, err, provider.ErrWeatherAPIStatusNotOk)
	assert.Nil(t, usecaseResponse)

	locationMock.AssertExpectations(t)
	temperatureMock.AssertExpectations(t)
}

func TestCelsiusToFahrenheit(t *testing.T) {
	cases := []struct {
		celsius  float64
		expected float64
	}{
		{0.0, 32.0},
		{25.0, 77.0},
		{100.0, 212.0},
	}
	for _, testCase := range cases {
		t.Run(fmt.Sprintf("should convert %.2f celsius to %.2f fahrenheit", testCase.celsius, testCase.expected), func(t *testing.T) {
			result := celsiusToFahrenheit(testCase.celsius)
			assert.Equal(t, testCase.expected, result)
		})
	}
}

func TestCelsiusToKelvin(t *testing.T) {
	cases := []struct {
		celsius  float64
		expected float64
	}{
		{0.0, 273.0},
		{25.0, 298.0},
		{100.0, 373.0},
	}
	for _, testCase := range cases {
		t.Run(fmt.Sprintf("should convert %.2f celsius to %.2f kelvin", testCase.celsius, testCase.expected), func(t *testing.T) {
			result := celsiusToKelvin(testCase.celsius)
			assert.Equal(t, testCase.expected, result)
		})
	}
}
