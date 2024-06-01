package configs

import (
	"log"

	"github.com/spf13/viper"
)

type cfg struct {
	WeatherAPIKey string
	WeatherAPIUrl string
	ViaCepAPIUrl  string
	ServerBPort   string
	ServerBHost   string
	ServerAPort   string
}

func GetConfig() *cfg {
	viper.SetDefault("WEATHER_API_URL", "https://api.weatherapi.com/v1/current.json")
	viper.SetDefault("VIA_CEP_API_URL", "https://viacep.com.br/ws/")
	viper.SetDefault("SERVER_A_PORT", "8080")
	viper.SetDefault("SERVER_B_PORT", "8081")
	viper.SetDefault("SERVER_B_HOST", "goserverb")

	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Println("Error reading config file", err)
	}

	return &cfg{
		WeatherAPIKey: viper.GetString("WEATHER_API_KEY"),
		WeatherAPIUrl: viper.GetString("WEATHER_API_URL"),
		ViaCepAPIUrl:  viper.GetString("VIA_CEP_API_URL"),
		ServerBPort:   viper.GetString("SERVER_B_PORT"),
		ServerBHost:   viper.GetString("SERVER_B_HOST"),
		ServerAPort:   viper.GetString("SERVER_A_PORT"),
	}
}
