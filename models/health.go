package models

type HealthStatus string

const (
	HealthStatusUp           HealthStatus = "UP"
	HealthStatusDown         HealthStatus = "DOWN"
	HealthStatusOutOfService HealthStatus = "OUT_OF_SERVICE"
	HealthStatusUnknown      HealthStatus = "UNKNOWN"
)

type HealthResponse struct {
	Status  HealthStatus   `json:"status"`
	Details map[string]any `json:"details"`
}
