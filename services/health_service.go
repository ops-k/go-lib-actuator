package services

import (
	"context"
	"slices"

	models "github.com/ops-k/go-lib-actuator/models"
)

type HealthIndicator interface {
	GetIdentifier(ctx context.Context) string
	GetHealthStatus(ctx context.Context) models.HealthResponse
}

type HealthIndicators []HealthIndicator
type LivenessHealthIndicators []HealthIndicator
type ReadinessHealthIndicators []HealthIndicator

// ActuatorService struct
type ActuatorHealthService struct {
	healthIndicators          HealthIndicators
	livenessHealthIndicators  LivenessHealthIndicators
	readinessHealthIndicators ReadinessHealthIndicators
}

func NewActuatorHealthService(
	healthIndicators HealthIndicators,
	livenessHealthIndicators LivenessHealthIndicators,
	readinessHealthIndicators ReadinessHealthIndicators,
) *ActuatorHealthService {
	return &ActuatorHealthService{
		healthIndicators:          healthIndicators,
		livenessHealthIndicators:  livenessHealthIndicators,
		readinessHealthIndicators: readinessHealthIndicators,
	}
}

func (svc *ActuatorHealthService) GetHealth(ctx context.Context) *models.HealthResponse {
	healthIndicatorResponses := make(map[string]models.HealthResponse)
	for _, healthIndicator := range svc.healthIndicators {
		healthIndicatorResponses[healthIndicator.GetIdentifier(ctx)] = healthIndicator.GetHealthStatus(ctx)
	}
	return mergeHealthResponses(healthIndicatorResponses, true)
}

func (svc *ActuatorHealthService) GetHealthLiveness(ctx context.Context) *models.HealthResponse {
	healthIndicatorResponses := make(map[string]models.HealthResponse)
	for _, healthIndicator := range svc.livenessHealthIndicators {
		healthIndicatorResponses[healthIndicator.GetIdentifier(ctx)] = healthIndicator.GetHealthStatus(ctx)
	}
	return mergeHealthResponses(healthIndicatorResponses, true)
}

func (svc *ActuatorHealthService) GetHealthReadiness(ctx context.Context) *models.HealthResponse {
	healthIndicatorResponses := make(map[string]models.HealthResponse)
	for _, healthIndicator := range svc.readinessHealthIndicators {
		healthIndicatorResponses[healthIndicator.GetIdentifier(ctx)] = healthIndicator.GetHealthStatus(ctx)
	}
	return mergeHealthResponses(healthIndicatorResponses, true)
}

var healthStatusPriorities = []models.HealthStatus{
	models.HealthStatusDown,
	models.HealthStatusOutOfService,
	models.HealthStatusUnknown,
	models.HealthStatusUp,
}

func mergeHealthResponses(healthIndicatorResponses map[string]models.HealthResponse, includeDetails bool) *models.HealthResponse {
	statusGroup := models.HealthStatusUp
	detailsGroup := make(map[string]interface{})
	for healthIndicatorIdentifier, healthIndicatorHealthResponse := range healthIndicatorResponses {
		if slices.Index(healthStatusPriorities, healthIndicatorHealthResponse.Status) < slices.Index(healthStatusPriorities, statusGroup) {
			statusGroup = healthIndicatorHealthResponse.Status
		}
		if includeDetails {
			detailsGroup[healthIndicatorIdentifier] = healthIndicatorHealthResponse.Details
		}
	}

	return &models.HealthResponse{
		Status:  statusGroup,
		Details: detailsGroup,
	}
}
