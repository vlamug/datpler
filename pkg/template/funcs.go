package template

import (
	"regexp"
	"strings"
	"text/template"
)

type extracter struct {
	compiledRegExpr map[string]*regexp.Regexp
}

func (e *extracter) extractRegExp(regExpr string, s string) string {
	if _, ok := e.compiledRegExpr[regExpr]; !ok {
		e.compiledRegExpr[regExpr] = regexp.MustCompile(regExpr)
	}

	return e.compiledRegExpr[regExpr].FindString(s)
}

func MakeTemplate(name string) *template.Template {
	extr := extracter{compiledRegExpr: make(map[string]*regexp.Regexp)}

	funcs := template.FuncMap{
		"contains":      strings.Contains,
		"extractRegExp": extr.extractRegExp,
	}

	return template.New(name).Funcs(funcs)
}
