package api

import (
	"encoding/json"
	"io"
	"net/http"
	"vkforwarder/internal/logger"
	"vkforwarder/internal/messaging"
)

type Options struct {
	Logger          logger.Logger
	MessProducer    *messaging.KafkaProducer
	VkMessTopicName string
	VkWebHookPath   string
	HealthCheckPath string
	Confirmation    string
	Secret          string
}

func SetupRoutes(opt Options) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc(opt.VkWebHookPath, func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			opt.Logger.Error("Failed to read request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var payload vkPayload

		if err := json.Unmarshal(body, &payload); err != nil {
			opt.Logger.Error("Failed to decode VK payload: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !verifyVKRequest(payload, opt.Secret, opt.Logger) {
			writeForbidden(w, opt.Logger)
			return
		}

		if isConfirmationRequest(payload) {
			writeConfirmation(w, opt.Logger, opt.Confirmation)
			return
		}

		if err := opt.MessProducer.Send(r.Context(), opt.VkMessTopicName, body); err != nil {
			opt.Logger.Error("Error delivering message to Kafka: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		opt.Logger.Info("Message delivered to Kafka")
		writeOK(w, opt.Logger)
	})

	mux.HandleFunc(opt.HealthCheckPath, func(w http.ResponseWriter, r *http.Request) {
		writeOK(w, opt.Logger)
	})

	return mux
}

type vkPayload struct {
	Type   string `json:"type"`
	Secret string `json:"secret"`
}

func isConfirmationRequest(payload vkPayload) bool {
	return payload.Type == "confirmation"
}

func verifyVKRequest(payload vkPayload, expectedSecret string, log logger.Logger) bool {
	if payload.Secret == expectedSecret {
		return true
	}

	log.Error("VK secret mismatch")

	return false
}

func writeConfirmation(w http.ResponseWriter, log logger.Logger, confirmation string) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(confirmation)); err != nil {
		log.Error("Failed to write confirmation response: %v", err)
	}
}

func writeOK(w http.ResponseWriter, log logger.Logger) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("ok")); err != nil {
		log.Error("Failed to write success response: %v", err)
	}
}

func writeForbidden(w http.ResponseWriter, log logger.Logger) {
	w.WriteHeader(http.StatusForbidden)
	if _, err := w.Write([]byte("forbidden")); err != nil {
		log.Error("Failed to write forbidden response: %v", err)
	}
}
