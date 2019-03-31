package metrics

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

type Storage struct {
	metrics map[string]*Metric
}

func NewStorage(metricSettings []*Metric) (*Storage, error) {
	builder := &Storage{metrics: make(map[string]*Metric)}

	err := builder.initMetrics(metricSettings)
	if err != nil {
		return nil, err
	}

	return builder, nil
}

func (s *Storage) initMetrics(metrics []*Metric) error {
	for _, metric := range metrics {
		var promMetric prometheus.Collector
		if metric.Type == counterMetricType {
			promMetric = prometheus.NewCounter(prometheus.CounterOpts{
				Name: metric.Name,
				Help: metric.Help,
			})
		} else if metric.Type == gaugeMetricType {
			promMetric = prometheus.NewGauge(prometheus.GaugeOpts{
				Name: metric.Name,
				Help: metric.Help,
			})
		} else if metric.Type == summaryMetricType {
			promMetric = prometheus.NewSummary(prometheus.SummaryOpts{
				Name:       metric.Name,
				Help:       metric.Help,
				Objectives: metric.ObjectivesValues(),
			})
		} else if metric.Type == histogramMetricType {
			promMetric = prometheus.NewHistogram(prometheus.HistogramOpts{
				Name:    metric.Name,
				Help:    metric.Help,
				Buckets: metric.Buckets,
			})
		} else if metric.Type == counterVecMetricType {
			promMetric = prometheus.NewCounterVec(prometheus.CounterOpts{
				Name: metric.Name,
				Help: metric.Help,
			}, metric.LabelNames())
		} else if metric.Type == gaugeVecMetricType {
			promMetric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Name: metric.Name,
				Help: metric.Help,
			}, metric.LabelNames())
		} else if metric.Type == summaryVecMetricType {
			promMetric = prometheus.NewSummaryVec(prometheus.SummaryOpts{
				Name:       metric.Name,
				Help:       metric.Help,
				Objectives: metric.ObjectivesValues(),
			}, metric.LabelNames())
		} else if metric.Type == histogramVecMetricType {
			promMetric = prometheus.NewHistogramVec(prometheus.HistogramOpts{
				Name:    metric.Name,
				Help:    metric.Help,
				Buckets: metric.Buckets,
			}, metric.LabelNames())
		} else {
			log.Printf("unsopperted metric type: %s\n", metric.Type)
			continue
		}

		// register metric
		prometheus.MustRegister(promMetric)

		metric.metric = promMetric

		// put metric into map
		s.metrics[metric.Name] = metric
	}

	return nil
}

func (s *Storage) Metrics() map[string]*Metric {
	return s.metrics
}
