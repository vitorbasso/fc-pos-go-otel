package webserver

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"otel/configs"
	"otel/internal/di"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewServerB() *http.Server {
	config := configs.GetConfig()
	temperatureHandler := di.NewTemperatureHandler(config.ViaCepAPIUrl, config.WeatherAPIUrl, config.WeatherAPIKey)
	r := chi.NewRouter()
	r.Use(middleware.Recoverer, middleware.Logger, middleware.Throttle(5), middleware.Heartbeat("/health-check"))
	r.Get("/temperatures/{cep}", temperatureHandler.GetTemperatureFromCep)
	return &http.Server{Addr: ":" + config.ServerBPort, Handler: r}
}

func NewServerA() *http.Server {
	config := configs.GetConfig()
	temperatureHandler := di.NewTemperatureHandler(config.ViaCepAPIUrl, config.WeatherAPIUrl, config.WeatherAPIKey)
	r := chi.NewRouter()
	r.Use(middleware.Recoverer, middleware.Logger, middleware.Throttle(5), middleware.Heartbeat("/health-check"))
	r.Post("/temperatures", temperatureHandler.GetTemperatureFromServerB)
	return &http.Server{Addr: ":" + config.ServerAPort, Handler: r}
}

func StartServer(server *http.Server) error {
	stop := make(chan os.Signal, 1)
	go func() {
		log.Printf("Server is running on port %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Panic(err)
		}
		log.Println("stopped receiving new requests")
	}()

	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("Shutting down server...")
	server.SetKeepAlivesEnabled(false)
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()
	if err := server.Shutdown(timeoutCtx); err != nil {
		return err
	}
	log.Println("Server gracefully stopped")
	return nil
}
