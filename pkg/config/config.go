package config

import (
	"log"

	"github.com/spf13/viper"
)

func Load() {
	viper.SetDefault("port", 8080)
	viper.SetDefault("lafzize_endpoint", "http://localhost:3001")
	viper.SetDefault("disable_csrf_checks", false)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("data")

	viper.SafeWriteConfig()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config: %v", err)
	}

	err = viper.WriteConfig()
	if err != nil {
		log.Fatalf("Error while writing config: %v", err)
	}
}
