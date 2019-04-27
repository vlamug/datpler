package metrics

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/vlamug/ratibor/pkg/template"

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
		var labelValues []string
		if needExpose, lvs := e.executeLabels(metric.Labels, data); !needExpose {
			continue
		} else {
			labelValues = lvs
		}

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

	tpl, err := template.MakeTemplate("metric").Parse(expr)
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

func (e *Processor) executeLabels(labels []Label, data map[string]string) (bool, []string) {
	labelValues := make([]string, len(labels))
	for k, lb := range labels {
		var (
			labelValue = lb.Value
			err        error
		)
		if len(lb.Values) != 0 {
			for _, val := range lb.Values {
				if val.Expr != "" {
					expr, err := e.executeExpr(val.Expr, data)
					if err != nil {
						log.Printf("could not execute one of label expr: %s\n", err)
						return false, nil
					}
					if expr == "" {
						continue
					}
				}

				labelValue, err = e.executeExpr(val.Value, data)
				if err != nil {
					log.Printf("could not execute one of label value: %s\n", err)
					return false, nil
				}

				break
			}

			// if there is no computed value, it means that the exposing should be skipped
			if labelValue == "" {
				log.Printf("there is no any value for label: %s\n", lb.Name)
				return false, nil
			}
		} else {

			if IsExecutable(lb.Value) {
				labelValue, err = e.executeExpr(lb.Value, data)
				if err != nil {
					log.Printf("could not execute label value: %s\n", err)
					return false, nil
				}
			}
		}

		labelValues[k] = labelValue
	}

	return true, labelValues
}
