package config

import (
	"fmt"
	"io/ioutil"
	"text/template"

	"bitbucket.org/plowdata/datpler/pkg/metrics"

	"gopkg.in/yaml.v2"
)

type Cfg struct {
	Input    Input             `yaml:"input"`
	Template Template          `yaml:"template"`
	Metrics  []*metrics.Metric `yaml:"metrics"`
}

func (cfg *Cfg) validate() error {
	if err := cfg.validateTemplate(); err != nil {
		return err
	}

	if err := cfg.validateMetrics(); err != nil {
		return err
	}

	return nil
}

func (cfg *Cfg) validateTemplate() error {
	_, err := template.New("template").Parse(cfg.Template.Pattern)

	return err
}

func (cfg *Cfg) validateMetrics() error {
	for _, metric := range cfg.Metrics {
		if err := metric.Validate(); err != nil {
			return fmt.Errorf("invalid metric sttings: %s, metric name: %s", err, metric.Name)
		}
	}

	return nil
}

type Input struct {
	Syslog []Syslog `yaml:"syslog"`
	API    []API    `yaml:"api"`
}

type Syslog struct {
	Name       string `yaml:"name"`
	ListenAddr string `yaml:"listenAddr"`
}

type API struct {
	Name       string `yaml:"name"`
	ListenAddr string `yaml:"listenAddr"`
	Path       string `yaml:"path"`
}

type Template struct {
	Type      string `yaml:"type"`
	Pattern   string `yaml:"pattern"`
	Delimiter string `yaml:"delimiter"`
}

func Load(path string) (*Cfg, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Cfg{}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	err = cfg.validate()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
