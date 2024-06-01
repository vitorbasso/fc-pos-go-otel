package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"otel/internal/gateway"
	"time"
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
	url := fmt.Sprintf("%s%s/json", v.viacepUrl, url.QueryEscape(cep))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("viacep_provider.FetchLocationByCep: %w", err)
	}
	client := http.Client{Timeout: 3 * time.Second}

	res, err := client.Do(req)
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
