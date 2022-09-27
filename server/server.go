package server

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"

	"github.com/chaosnative/toxiproxy/v2"
	"github.com/chaosnative/toxiproxy/v2/collectors"
)

type cliArguments struct {
	host           string
	port           string
	config         string
	seed           int64
	printVersion   bool
	proxyMetrics   bool
	runtimeMetrics bool
}

func parseArguments(host, port string) cliArguments {
	return cliArguments{
		host:           host,
		port:           port,
		config:         "",
		seed:           time.Now().UTC().UnixNano(),
		printVersion:   false,
		proxyMetrics:   false,
		runtimeMetrics: false,
	}
}

func ProxyServer(exitSignal chan os.Signal, host, port string) {
	// Handle SIGTERM to exit cleanly
	go func() {
		for {
			select {
			case <-exitSignal:
				return
			}
		}

	}()

	cli := parseArguments(host, port)
	run(cli)
}

func run(cli cliArguments) {
	if cli.printVersion {
		fmt.Printf("toxiproxy-server version %s\n", toxiproxy.Version)
		return
	}

	setupLogger()

	rand.Seed(cli.seed)

	metrics := toxiproxy.NewMetricsContainer(prometheus.NewRegistry())
	server := toxiproxy.NewServer(metrics)
	if cli.proxyMetrics {
		server.Metrics.ProxyMetrics = collectors.NewProxyMetricCollectors()
	}
	if cli.runtimeMetrics {
		server.Metrics.RuntimeMetrics = collectors.NewRuntimeMetricCollectors()
	}
	if len(cli.config) > 0 {
		server.PopulateConfig(cli.config)
	}

	server.Listen(cli.host, cli.port)
}

func setupLogger() {
	val, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		return
	}

	lvl, err := logrus.ParseLevel(val)
	if err == nil {
		logrus.SetLevel(lvl)
		return
	}

	valid_levels := make([]string, len(logrus.AllLevels))
	for i, level := range logrus.AllLevels {
		valid_levels[i] = level.String()
	}
	levels := strings.Join(valid_levels, ",")

	logrus.Errorf("unknown LOG_LEVEL value: \"%s\", use one of: %s", val, levels)
}
