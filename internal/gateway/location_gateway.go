package gateway

import (
	"context"
	"errors"
)

var (
	ErrCepNotFound = errors.New("cep not found")
)

type LocationResponse struct {
	Cep          string
	State        string
	City         string
	Neighborhood string
	Street       string
}

type LocationGateway interface {
	FetchLocationByCep(ctx context.Context, cep string) (*LocationResponse, error)
}
