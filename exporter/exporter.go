package exporter

import (
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

type Exporter struct {
	up *prometheus.Desc
}

func New(apiKey string, logger log.Logger) *Exporter {
	return &Exporter{}
}

func (e *Exporter) Describe(chan<- *prometheus.Desc) {

}

func (e *Exporter) Collect(chan<- prometheus.Metric) {

}
