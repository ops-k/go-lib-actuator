package metrics

import (
	"time"

	services "github.com/ops-k/go-lib-actuator/services"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
)

type UptimeMetricsCollector struct {
	logger    zerolog.Logger
	startTime time.Time
	desc      *prometheus.Desc
}

var (
	uptimeDesc = prometheus.NewDesc(
		"uptime_seconds",
		"Application uptime in seconds",
		nil,
		nil,
	)
)

func NewUptimeMetricsCollector(
	logger zerolog.Logger,
	actuatorHealthService *services.ActuatorHealthService,
) *UptimeMetricsCollector {
	return &UptimeMetricsCollector{
		logger:    logger,
		startTime: time.Now(),
		desc:      uptimeDesc,
	}
}

// Describe sends the metric description to Prometheus.
func (c *UptimeMetricsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.desc
}

// Collect is called by Prometheus to collect the metric value.
func (c *UptimeMetricsCollector) Collect(ch chan<- prometheus.Metric) {
	uptime := time.Since(c.startTime).Seconds()
	ch <- prometheus.MustNewConstMetric(
		c.desc,
		prometheus.GaugeValue,
		uptime,
	)
}
