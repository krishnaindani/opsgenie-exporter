package main

import (
	"fmt"
	"github.com/go-kit/log/level"
	"github.com/krishnaindani/opsgenie-exporter/exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	webflag "github.com/prometheus/exporter-toolkit/web/kingpinflag"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"os"
)

func main() {

	var (
		listenAddress = kingpin.Flag("web.listen-address", "Address to listen for web and telemetry").Default(":9201").String()
		metricsPath   = kingpin.Flag("web.metrics-path", "Path to expose metrics on").Default("/metrics").String()
		apiKey        = kingpin.Flag("api-key", "Opsgenie api key").Default("").String()
		webConfig     = webflag.AddFlags(kingpin.CommandLine)
	)

	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.HelpFlag.Short('h')
	kingpin.Version(version.Print("opsgenie_exporter"))
	kingpin.Parse()

	logger := promlog.New(promlogConfig)

	level.Info(logger).Log("msg", "Starting opsgenie_exporter", "version", version.Info())
	level.Info(logger).Log("msg", "Build context", "context", version.BuildContext())

	prometheus.MustRegister(version.NewCollector("opsgenie_exporter"))
	exp, err := exporter.New(*apiKey, logger)
	if err != nil {
		level.Error(logger).Log("msg", "Error creating the exporter", "err", err)
		os.Exit(1)
	}
	prometheus.MustRegister(exp)

	http.Handle(*metricsPath, promhttp.Handler())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Opsgenie Exporter</title></head>
             <body>
             <h1>Opsgenie Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})

	http.HandleFunc("/-/healthy", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Healthy")
	})

	http.HandleFunc("/-/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Ready")
	})

	level.Info(logger).Log("msg", "Listening on address", "address", *listenAddress)
	server := &http.Server{
		Addr: *listenAddress,
	}

	if err := web.ListenAndServe(server, *webConfig, logger); err != nil {
		level.Error(logger).Log("msg", "Error running http server", err)
		os.Exit(1)
	}
}
