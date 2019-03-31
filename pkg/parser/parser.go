package parser

type Parser interface {
	Parse(content string) (map[string]string, error)
}
