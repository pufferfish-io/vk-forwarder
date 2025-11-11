package main

import (
	"context"
	"log"
	"net/http"

	"vkforwarder/internal/api"
	"vkforwarder/internal/config"
	"vkforwarder/internal/logger"
	"vkforwarder/internal/messaging"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	logger, clean := logger.NewZapLogger()
	defer clean()

	opt := messaging.Option{
		Context:      context.Background(),
		Logger:       logger,
		Broker:       cfg.Kafka.BootstrapServersValue,
		SaslUsername: cfg.Kafka.SaslUsername,
		SaslPassword: cfg.Kafka.SaslPassword,
	}

	prod, err := messaging.NewKafkaProducer(opt)

	if err != nil {
		log.Fatalf("config: %v", err)
	}
	defer prod.Close()

	mux := api.SetupRoutes(api.Options{
		Logger:          logger,
		MessProducer:    prod,
		VkMessTopicName: cfg.Kafka.VkMessTopicName,
		VkWebHookPath:   cfg.Api.VkWebHookPath,
		HealthCheckPath: cfg.Api.HealthCheckPath,
		Confirmation:    cfg.VK.Confirmation,
		Secret:          cfg.VK.Secret,
	})
	log.Printf("🌐 Webhook server is listening on %s...", cfg.Server.Addr)
	log.Fatal(http.ListenAndServe(cfg.Server.Addr, mux))
}
