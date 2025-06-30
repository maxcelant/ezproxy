package proxy

type Logger interface {
	Info(string, ...any)
	Error(string, ...any)
}
