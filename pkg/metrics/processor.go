package metrics

import (
	"bytes"
	"fmt"
	"github.com/prometheus/common/log"
	"io/ioutil"
	"strconv"
	"text/template"

	"github.com/prometheus/client_golang/prometheus"
)

// Processor evaluates expr+value and exposes metrics
type Processor struct {
	storage *Storage
}

func NewProcessor(storage *Storage) *Processor {
	return &Processor{storage: storage}
}

func (e *Processor) Process(data map[string]string) error {
	for _, metric := range e.storage.Metrics() {
		executedExpr, err := e.executeExpr(metric.Expr, data)
		if err != nil {
			return fmt.Errorf("could not execute metric expression: %s", err)
		}

		executedValue, err := e.executeExpr(metric.Value, data)
		if err != nil {
			return fmt.Errorf("could not execute metric value: %s", err)
		}

		needExpose := false
		exposeValue := float64(0)
		if metric.Expr != "" && metric.Value != "" {
			// if expr is true, then expose metric
			if executedExpr != "" {
				needExpose = true
				val, err := strconv.ParseFloat(executedValue, 64)
				// if it is possible to parse value
				if err == nil {
					exposeValue = val
				} else {
					// else check if it is true
					if executedValue != "" {
						exposeValue = float64(1)
					} else {
						exposeValue = float64(0)
					}
				}
			}
		} else if metric.Expr != "" {
			// if it is not empty
			if executedExpr != "" {
				needExpose = true
				val, err := strconv.ParseFloat(executedExpr, 64)
				// if it is possible to parse value
				if err == nil {
					exposeValue = val
				} else {
					exposeValue = float64(1)
				}
			}
		} else if metric.Value != "" {
			needExpose = true
			// if it is not empty
			if executedValue != "" {
				val, err := strconv.ParseFloat(executedValue, 64)
				// if it is possible to parse value
				if err == nil {
					exposeValue = val
				} else {
					exposeValue = float64(1)
				}
			} else {
				exposeValue = float64(0)
			}
		}

		if !needExpose {
			continue
		}

		// if it is not vector type, empty slice will be returned
		labelValues := e.executeLabels(metric.Labels, data)

		switch metric.Type {
		case counterMetricType:
			metric.metric.(prometheus.Counter).Add(exposeValue)
		case gaugeMetricType:
			metric.metric.(prometheus.Gauge).Set(exposeValue)
		case summaryMetricType:
			metric.metric.(prometheus.Summary).Observe(exposeValue)
		case histogramMetricType:
			metric.metric.(prometheus.Histogram).Observe(exposeValue)
		case counterVecMetricType:
			metric.metric.(*prometheus.CounterVec).WithLabelValues(labelValues...).Add(exposeValue)
		case gaugeVecMetricType:
			metric.metric.(*prometheus.GaugeVec).WithLabelValues(labelValues...).Set(exposeValue)
		case summaryVecMetricType:
			metric.metric.(*prometheus.SummaryVec).WithLabelValues(labelValues...).Observe(exposeValue)
		case histogramVecMetricType:
			metric.metric.(*prometheus.HistogramVec).WithLabelValues(labelValues...).Observe(exposeValue)
		}
	}

	return nil
}

// executeExpr tries to execute passes expression with variables and values in data
func (e *Processor) executeExpr(expr string, data map[string]string) (string, error) {
	if expr == "" {
		return "", nil
	}

	if !IsExecutable(expr) {
		return expr, nil
	}

	tpl, err := template.New("metric").Parse(expr)
	if err != nil {
		return "", fmt.Errorf("could not parse expression: %s, error: %s\n", expr, err)
	}

	buf := &bytes.Buffer{}
	err = tpl.Execute(buf, data)
	if err != nil {
		return "", fmt.Errorf("could not execute expression: %s, error: %s\n", expr, err)
	}

	res, err := ioutil.ReadAll(buf)
	if err != nil {
		return "", fmt.Errorf("could not read executed expression: %s, error: %s", expr, err)
	}

	return string(res), nil
}

func (e *Processor) executeLabels(labels []Label, data map[string]string) []string {
	labelValues := make([]string, len(labels))
	for k, lb := range labels {
		var (
			err        error
			labelValue = lb.Value
		)

		if IsExecutable(lb.Value) {
			labelValue, err = e.executeExpr(lb.Value, data)
			if err != nil {
				log.Warnf("could not execute label value: %s", labelValue)
			}
		}

		labelValues[k] = labelValue
	}

	return labelValues
}
