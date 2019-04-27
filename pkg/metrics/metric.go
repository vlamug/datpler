package metrics

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/vlamug/ratibor/pkg/template"
)

var validMetricNameRegexp = regexp.MustCompile(`^[\w_]+$`)

type Metric struct {
	Name       string      `yaml:"name"`
	Type       string      `yaml:"type"`
	Help       string      `yaml:"help"`
	Expr       string      `yaml:"expr"`
	Value      string      `yaml:"value"`
	Objectives []Objective `yaml:"objectives"`
	Buckets    []float64   `yaml:"buckets"`
	Labels     []Label     `yaml:"labels"`

	metric interface{}
}

type Objective struct {
	Quantile float64 `yaml:"quantile"`
	Error    float64 `yaml:"error"`
}

type Label struct {
	Name   string  `yaml:"name"`
	Value  string  `yaml:"value"`
	Values []Value `yaml:"values"`
}

type Value struct {
	Value string `yaml:"value"`
	Expr  string `yaml:"expr"`
}

func (m *Metric) ObjectivesValues() map[float64]float64 {
	objectives := make(map[float64]float64)
	for _, objective := range m.Objectives {
		objectives[objective.Quantile] = objective.Error
	}

	return objectives
}

func (m *Metric) LabelNames() []string {
	names := make([]string, len(m.Labels))
	for k, lb := range m.Labels {
		names[k] = lb.Name
	}

	return names
}

func (m *Metric) Validate() error {
	if !validateType(m.Type) {
		return fmt.Errorf("invalid metric type: %s", m.Type)
	}

	if !validateName(m.Name) {
		return fmt.Errorf("invalid metric name: %s, alphabetic symbols and underscore are allowed", m.Name)
	}

	if err := validateExecutableValue(m.Expr); err != nil {
		return fmt.Errorf("invalid expression: %s", err)
	}

	if err := validateExecutableValue(m.Value); err != nil {
		return fmt.Errorf("invalid value: %s", m.Value)
	}

	// validate metric type specific settings
	if m.Type == summaryMetricType {
		if err := validateObjectives(m.Objectives); err != nil {
			return err
		}
	} else if m.Type == histogramMetricType {
		if err := validateBuckets(m.Buckets); err != nil {
			return err
		}
	}

	if IsVectorType(m.Type) {
		if err := m.validateLabels(); err != nil {
			return err
		}
	} else {
		if len(m.Labels) > 0 {
			return errors.New("specified labels are redundant for not vector metric")
		}
	}

	return nil
}

func (m *Metric) validateLabels() error {
	for _, lb := range m.Labels {
		if err := validateExecutableValue(lb.Value); err != nil {
			return errors.New(fmt.Sprintf("label value cannot be executable: %s", err))
		}

		for _, val := range lb.Values {
			if val.Value == "" && val.Expr == "" {
				return errors.New("`value` or `expr` in label value are required")
			}
			if err := validateExecutableValue(val.Value); err != nil {
				return errors.New(fmt.Sprintf("one of the label value cannot be executable: %s", err))
			}
			if err := validateExecutableValue(val.Expr); err != nil {
				return errors.New(fmt.Sprintf("one of the label expr cannot be executable: %s", err))
			}
		}
	}

	return nil
}

// @todo make validate functions as method of struct
func validateObjectives(objectives []Objective) error {
	if len(objectives) == 0 {
		return errors.New("objectives are not specified")
	}

	return nil
}

func validateBuckets(buckets []float64) error {
	if len(buckets) == 0 {
		return errors.New("buckets are not specified")
	}

	return nil
}

func validateType(mType string) bool {
	for _, t := range AllMetricTypes {
		if mType == t {
			return true
		}
	}

	return false
}

func validateName(name string) bool {
	return validMetricNameRegexp.MatchString(name)
}

func validateExecutableValue(value string) error {
	if value != "" && IsExecutable(value) {
		_, err := template.MakeTemplate("value").Parse(value)
		if err != nil {
			return err
		}
	}

	return nil
}
