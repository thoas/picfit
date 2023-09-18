package logger

type Config struct {
	Level       string `json:"level"`
	ContextKeys []string
}
