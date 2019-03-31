package data

import (
	"bitbucket.org/plowdata/datpler/config"
	"bitbucket.org/plowdata/datpler/pkg/input"
	"bitbucket.org/plowdata/datpler/pkg/metrics"
	"bitbucket.org/plowdata/datpler/pkg/parser"

	"github.com/prometheus/common/log"
)

type Plower struct {
	inCfg config.Input

	psr       parser.Parser
	processor *metrics.Processor

	syslogInput func(listenAddr string, logLines chan<- *input.LogLine)
	APIInput    func(listenAddr, path string, lines chan<- string)
}

func NewPlower(
	inCfg config.Input,
	psr parser.Parser,
	processor *metrics.Processor,
	syslogInput func(listenAddr string, logLines chan<- *input.LogLine),
	APIInput func(listenAddr, path string, lines chan<- string),
) *Plower {
	return &Plower{inCfg: inCfg, psr: psr, processor: processor, syslogInput: syslogInput, APIInput: APIInput}
}

func (h *Plower) Plow() {
	if len(h.inCfg.Syslog) > 0 {
		for _, in := range h.inCfg.Syslog {
			go func(name, listenAddr string) {
				logLines := make(chan *input.LogLine)
				go h.syslogInput(listenAddr, logLines)

				for line := range logLines {
					err := h.plowData(line.Data)
					if err != nil {
						log.Errorf("error occurred during plow data for syslog input: %s, error: %s", name, err)
					}
				}
			}(in.Name, in.ListenAddr)
		}
	}

	if len(h.inCfg.API) > 0 {
		for _, in := range h.inCfg.API {
			go func(name, listenAddr, path string) {
				lines := make(chan string)

				go h.APIInput(listenAddr, path, lines)

				for line := range lines {
					err := h.plowData(line)
					if err != nil {
						log.Errorf("error occurred during plow data for api input: %s", name)
					}
				}
			}(in.Name, in.ListenAddr, in.Path)
		}
	}
}

func (h *Plower) plowData(data string) error {
	res, err := h.psr.Parse(data)
	if err != nil {
		return err
	}

	return h.processor.Process(res)
}
