package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	Addr string `validate:"required" env:"SERVER_ADDR_VK_FORWARDER" envDefault:":8080"`
}

type Kafka struct {
	BootstrapServersValue string `validate:"required" env:"BOOTSTRAP_SERVERS_VALUE"`
	VkMessTopicName       string `validate:"required" env:"TOPIC_NAME_VK_UPDATES"`
	SaslUsername          string `validate:"required" env:"SASL_USERNAME"`
	SaslPassword          string `validate:"required" env:"SASL_PASSWORD"`
}

type VK struct {
	Confirmation string `validate:"required" env:"CONFIRMATION"`
	Secret       string `validate:"required" env:"SECRET"`
}

type Api struct {
	VkWebHookPath   string `validate:"required" env:"PATH_VK_FORWARDER_VK_WEB_HOOK"`
	HealthCheckPath string `validate:"required" env:"PATH_VK_FORWARDER_HEALTH_CHECK"`
}
type Config struct {
	Server Server
	Kafka  Kafka `envPrefix:"KAFKA_"`
	VK     VK    `envPrefix:"VK_"`
	Api    Api   `envPrefix:"API_"`
}

func Load() (*Config, error) {
	var c Config
	if err := env.Parse(&c); err != nil {
		return nil, fmt.Errorf("env parse: %w", err)
	}

	v := validator.New()
	if err := v.Struct(c); err != nil {
		return nil, fmt.Errorf("config validate: %w", err)
	}

	return &c, nil
}
