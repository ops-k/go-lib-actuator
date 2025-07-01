package metrics

import (
	"context"
	"sync"

	models "github.com/ops-k/go-lib-actuator/models"
	services "github.com/ops-k/go-lib-actuator/services"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
)

type HealthStatusProvider struct {
	prometheusDesc         *prometheus.Desc
	healthStatusProviderFn func(context.Context) *models.HealthResponse
}

type HealthMetricsCollector struct {
	logger                zerolog.Logger
	actuatorHealthService *services.ActuatorHealthService
	healthStatusProviders []HealthStatusProvider
}

var (
	upDesc = prometheus.NewDesc(
		"up",
		"Health status of the application",
		nil,
		nil,
	)

	aliveDesc = prometheus.NewDesc(
		"alive",
		"Liveness status of the application",
		nil,
		nil,
	)

	readyDesc = prometheus.NewDesc(
		"ready",
		"Readiness status of the application",
		nil,
		nil,
	)
)

func NewHealthMetricsCollector(
	logger zerolog.Logger,
	actuatorHealthService *services.ActuatorHealthService,
) *HealthMetricsCollector {
	// list all providers and their related metric
	return &HealthMetricsCollector{
		logger:                logger,
		actuatorHealthService: actuatorHealthService,
		healthStatusProviders: []HealthStatusProvider{
			{
				prometheusDesc:         upDesc,
				healthStatusProviderFn: actuatorHealthService.GetHealth,
			},
			{
				prometheusDesc:         aliveDesc,
				healthStatusProviderFn: actuatorHealthService.GetHealthLiveness,
			},
			{
				prometheusDesc:         readyDesc,
				healthStatusProviderFn: actuatorHealthService.GetHealthReadiness,
			},
		},
	}
}

// Describe sends the metric description to Prometheus.
func (c *HealthMetricsCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, provider := range c.healthStatusProviders {
		ch <- provider.prometheusDesc
	}
}

// Collect is called by Prometheus to collect the metric value.
func (c *HealthMetricsCollector) Collect(ch chan<- prometheus.Metric) {
	wg := sync.WaitGroup{}
	wg.Add(len(c.healthStatusProviders))

	for _, provider := range c.healthStatusProviders {
		go func(prometheusDesc *prometheus.Desc, healthStatusProvider func(context.Context) *models.HealthResponse) {
			defer wg.Done()
			healthResponse := healthStatusProvider(context.TODO())
			value := mapHealthStatusToValue(healthResponse.Status)
			ch <- prometheus.MustNewConstMetric(
				prometheusDesc,
				prometheus.GaugeValue,
				value,
			)
		}(provider.prometheusDesc, provider.healthStatusProviderFn)
	}

	wg.Wait()
}

func mapHealthStatusToValue(status models.HealthStatus) float64 {
	switch status {
	case models.HealthStatusUp:
		return 1
	case models.HealthStatusDown:
		return 0
	case models.HealthStatusUnknown:
		return -1
	default:
		return -1
	}
}
