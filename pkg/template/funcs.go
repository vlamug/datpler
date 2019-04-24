package template

import (
	"regexp"
	"strings"
	"text/template"
)

type extractor struct {
	compiledRegExpr map[string]*regexp.Regexp
}

func (e *extractor) extractRegExp(regExpr string, s string) string {
	if _, ok := e.compiledRegExpr[regExpr]; !ok {
		e.compiledRegExpr[regExpr] = regexp.MustCompile(regExpr)
	}

	return e.compiledRegExpr[regExpr].FindString(s)
}

// MakeTemplate makes template extended with functions
func MakeTemplate(name string) *template.Template {
	extr := extractor{compiledRegExpr: make(map[string]*regexp.Regexp)}

	funcs := template.FuncMap{
		"contains":      strings.Contains,
		"extractRegExp": extr.extractRegExp,
	}

	return template.New(name).Funcs(funcs)
}
