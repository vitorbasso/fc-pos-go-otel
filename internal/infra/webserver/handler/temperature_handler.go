package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"otel/internal/usecase"
	"regexp"

	"github.com/go-chi/chi/v5"
)

type TemperatureHandlerInput struct {
	Cep string `json:"cep"`
}

type TemperateHandlerResponse struct {
	City       string `json:"city"`
	Celsius    string `json:"temp_C"`
	Fahrenheit string `json:"temp_F"`
	Kelvin     string `json:"temp_K"`
}

var (
	ErrInvalidZipCode = errors.New("invalid zipcode")
)

type TemperatureHandler struct {
	getTemperatureFromCepUseCase     usecase.GetTemperatureFromCepUseCase
	getTemperatureFromServerBUseCase *usecase.GetTemperatureFromServerBUseCase
}

func NewTemperatureHandler(getTemperatureFromCepUseCase usecase.GetTemperatureFromCepUseCase, getTemperatureFromServerBUseCase *usecase.GetTemperatureFromServerBUseCase) *TemperatureHandler {
	return &TemperatureHandler{
		getTemperatureFromCepUseCase:     getTemperatureFromCepUseCase,
		getTemperatureFromServerBUseCase: getTemperatureFromServerBUseCase,
	}
}

func (t *TemperatureHandler) GetTemperatureFromCep(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	cep := chi.URLParam(r, "cep")
	if cep == "" || !isValidCep(cep) {
		errString := ErrInvalidZipCode.Error()
		log.Printf("error: %s", errString)
		http.Error(w, errString, http.StatusUnprocessableEntity)
		return
	}
	response, err := t.getTemperatureFromCepUseCase.Execute(r.Context(), cep)
	if errors.Is(err, usecase.ErrNotFound) {
		log.Printf("error: %s", err.Error())
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("error: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(TemperateHandlerResponse{
		City:       response.City,
		Celsius:    response.Celsius,
		Fahrenheit: response.Fahrenheit,
		Kelvin:     response.Kelvin,
	}); err != nil {
		log.Printf("error: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (t *TemperatureHandler) GetTemperatureFromServerB(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var input TemperatureHandlerInput
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Println("Error decoding input", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if input.Cep == "" || !isValidCep(input.Cep) {
		errString := ErrInvalidZipCode.Error()
		log.Printf("error: %s", errString)
		http.Error(w, errString, http.StatusUnprocessableEntity)
		return
	}
	response, err := t.getTemperatureFromServerBUseCase.Execute(r.Context(), input.Cep)
	if err != nil {
		log.Printf("error: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = json.NewEncoder(w).Encode(TemperateHandlerResponse{
		City:       response.City,
		Celsius:    response.Celsius,
		Fahrenheit: response.Fahrenheit,
		Kelvin:     response.Kelvin,
	}); err != nil {
		log.Printf("error: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

const validCepRegex string = `^\d{5}-?\d{3}$`

var (
	cepRegex = regexp.MustCompile(validCepRegex)
)

func isValidCep(cep string) bool {
	return cepRegex.MatchString(cep)
}
