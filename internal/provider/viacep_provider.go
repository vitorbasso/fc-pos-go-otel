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
	ErrViaCepStatusNotOk = errors.New("viacep status not ok")
	ErrViaCepNotFound    = errors.New("viacep not found")
)

type viacepResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
	Erro        bool   `json:"erro"`
}

type ViaCepLocationProvider struct {
	viacepUrl string
}

func NewViaCepLocationProvider(viacepUrl string) *ViaCepLocationProvider {
	return &ViaCepLocationProvider{
		viacepUrl: viacepUrl,
	}
}

func (v *ViaCepLocationProvider) FetchLocationByCep(ctx context.Context, cep string) (*gateway.LocationResponse, error) {
	tracer := ctx.Value(constants.CtxTracerKey).(trace.Tracer)
	ctx, span := tracer.Start(ctx, "viacep_provider.FetchLocationByCep")
	defer span.End()

	url := fmt.Sprintf("%s%s/json", v.viacepUrl, url.QueryEscape(cep))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("viacep_provider.FetchLocationByCep: %w", err)
	}
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("viacep_provider.FetchLocationByCep: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("viacep_provider.FetchLocationByCep: %w; code: %s", ErrViaCepStatusNotOk, res.Status)
	}
	var response viacepResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("viacep_provider.FetchLocationByCep: %w", err)
	}
	if response.Erro {
		return nil, fmt.Errorf("viacep_provider.FetchLocationByCep: %w", gateway.ErrCepNotFound)
	}
	return &gateway.LocationResponse{
		Cep:          response.Cep,
		State:        response.Uf,
		City:         response.Localidade,
		Neighborhood: response.Bairro,
		Street:       response.Logradouro,
	}, nil
}
