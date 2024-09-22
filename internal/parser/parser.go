package parser

type Parser interface {
	Parse(data []byte, taskId string) ([]Finding, error)
}
