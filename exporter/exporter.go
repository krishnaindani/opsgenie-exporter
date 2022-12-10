package exporter

import (
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/krishnaindani/opsgenie-exporter/client/teams"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace      = "opsgenie_exporter"
	teamsSubSystem = "teams"
)

type Exporter struct {
	teamClient *teams.Client
	logger     log.Logger

	up           *prometheus.Desc
	countOfTeams *prometheus.Desc
}

// New returns initialized exporter
func New(apiKey string, logger log.Logger) (*Exporter, error) {

	teamClient, err := teams.NewClient(teams.Config{
		ApiKey: apiKey,
	})

	if err != nil {
		level.Error(logger).Log("msg", "Error creating team client", "err", err)
		return nil, err
	}

	return &Exporter{
		teamClient: teamClient,

		up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "up"),
			"Successful getting all metrics from opsgenie",
			nil,
			nil,
		),

		countOfTeams: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, teamsSubSystem, "count"),
			"Count of all teams in opsgenie",
			nil,
			nil,
		),
	}, nil
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.up
	ch <- e.countOfTeams
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	ok := e.collectTotalTeamsMetric(ch)

	if ok {
		ch <- prometheus.MustNewConstMetric(
			e.up,
			prometheus.GaugeValue,
			1.0,
		)
	} else {
		ch <- prometheus.MustNewConstMetric(
			e.up,
			prometheus.GaugeValue,
			0.0,
		)
	}
}

func (e *Exporter) collectTotalTeamsMetric(ch chan<- prometheus.Metric) bool {
	teams, err := e.teamClient.GetCountOfAllTeams()
	if err != nil {
		level.Error(e.logger).Log("msg", "Error getting count of total teams metric", "err", err)
		return false
	}

	ch <- prometheus.MustNewConstMetric(
		e.countOfTeams,
		prometheus.GaugeValue,
		teams,
	)

	return true
}
