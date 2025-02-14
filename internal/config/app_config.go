package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	APIURL            string
	MaxRetries        int
	RetryDelaySeconds int
	RetryDelay        time.Duration
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("config")

	// Set default values.
	viper.SetDefault("maxRetries", 3)
	viper.SetDefault("retryDelaySeconds", 1)

	// Read the config file.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Warning: error reading config file: %v\n", err)
	}

	// Read env file
	viper.AutomaticEnv()
	viper.BindEnv("api.url", "API_URL")

	cfg := &Config{
		APIURL:            viper.GetString("api.url"),
		MaxRetries:        viper.GetInt("maxRetries"),
		RetryDelaySeconds: viper.GetInt("retryDelaySeconds"),
	}
	cfg.RetryDelay = time.Duration(cfg.RetryDelaySeconds) * time.Second

	if cfg.APIURL == "" {
		return nil, fmt.Errorf("API URL must be set in environment variable API_URL or in the config file")
	}

	return cfg, nil
}
