package metrics

import (
	models "github.com/ops-k/go-lib-actuator/models"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
)

type BuildInfoMetricsCollector struct {
	logger    zerolog.Logger
	buildInfo *models.BuildInfo
	desc      *prometheus.Desc
}

var (
	buildInfoDesc = prometheus.NewDesc(
		"build_info",
		"Build information",
		[]string{"name", "version", "commit", "date"},
		nil,
	)
)

func NewBuildInfoMetricsCollector(
	logger zerolog.Logger,
	buildInfo *models.BuildInfo,
) *BuildInfoMetricsCollector {
	return &BuildInfoMetricsCollector{
		logger:    logger,
		buildInfo: buildInfo,
		desc:      buildInfoDesc,
	}
}

// Describe sends the metric description to Prometheus.
func (c *BuildInfoMetricsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.desc
}

// Collect is called by Prometheus to collect the metric value.
func (c *BuildInfoMetricsCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		c.desc,
		prometheus.GaugeValue,
		float64(1),
		c.buildInfo.Name,
		c.buildInfo.Version,
		c.buildInfo.Commit,
		c.buildInfo.Date,
	)
}
