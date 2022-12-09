package log

type Logger interface {
	Debug(msg string)
	Info(msg string, fields ...Field)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
}

type Field struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
