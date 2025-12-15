package internal

import (
	"encoding/json"
	"time"
)

// HelloHandler is an example HTTP handler for demonstration purposes.
func HelloHandler() string {
	return "Hello from internal package!"
}

// HealthCheckResponse represents the health check response structure.
type HealthCheckResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
	Version   string    `json:"version"`
}

// HealthCheckHandler returns the health status of the service.
func HealthCheckHandler() ([]byte, error) {
	response := HealthCheckResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Service:   "Go Microservice",
		Version:   "1.0.0",
	}
	return json.Marshal(response)
}
