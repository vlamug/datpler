package parser

import (
	"log"
	"strings"
)

type Plain struct {
	tpl       string
	delimiter string
	vars      []string
}

func NewPlain(tpl, delimiter string) *Plain {
	vars := strings.Split(tpl, delimiter)
	for k, v := range vars {
		vars[k] = strings.Trim(strings.Trim(strings.Trim(strings.Trim(v, " "), "{"), "}"), ".")
	}

	return &Plain{tpl: tpl, delimiter: delimiter, vars: vars}
}

func (e *Plain) Parse(data string) (map[string]string, error) {
	vals := strings.Split(data, e.delimiter)

	if len(e.vars) != len(vals) {
		log.Printf("dismatch var names and values")
		return nil, nil
	}

	varsVals := make(map[string]string)
	for k, v := range e.vars {
		varsVals[v] = strings.Trim(vals[k], " ")
	}

	return varsVals, nil
}
