package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/vlamug/ratibor/config"
	"github.com/vlamug/ratibor/pkg/data"
	"github.com/vlamug/ratibor/pkg/input"
	"github.com/vlamug/ratibor/pkg/metrics"
	"github.com/vlamug/ratibor/pkg/parser"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	defaultListenAddr = ":11110"
	defaultConfigPath = "etc/config.yml"
)

func main() {
	var (
		listenAddr string
		configPath string
	)

	flag.StringVar(&listenAddr, "listen.addr", defaultListenAddr, "Address to listen requests")
	flag.StringVar(&configPath, "config.path", defaultConfigPath, "Config file path")
	flag.Parse()

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalln(err)
	}

	parserFactory := parser.NewFactory()
	parserFactory.AddParser(parser.PlainTplType, func() parser.Parser {
		return parser.NewPlain(cfg.Template.Pattern, cfg.Template.Delimiter)
	})
	psr, err := parserFactory.Create(cfg.Template.Type)
	if err != nil {
		log.Fatalf("could not create parser: %s", err)
	}

	// init metrics storage
	storage, err := metrics.NewStorage(cfg.Metrics)
	if err != nil {
		log.Fatalln(err)
	}

	// inputs
	syslogInput := input.NewSyslog()
	apiInput := input.NewAPI()

	// init and run data plower
	plower := data.NewPlower(cfg.Input, psr, metrics.NewProcessor(storage), syslogInput, apiInput)
	plower.Plow()

	http.Handle("/metrics", promhttp.Handler())

	log.Printf("Listening web on: %s", listenAddr)
	log.Fatalln(http.ListenAndServe(listenAddr, nil))
}
