package logger

import "context"

// Logger provides a leveled-logging interface.
type Logger interface {
	// standard logger methods
	Print(args ...interface{})
	Printf(format string, args ...interface{})
	Println(args ...interface{})

	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatalln(args ...interface{})

	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Panicln(args ...interface{})

	// Leveled methods, from logrus
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Debugln(args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Errorln(args ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Infoln(args ...interface{})

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Warnln(args ...interface{})
}

const key = "logger"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the KVStore associated with this context.
func FromContext(c context.Context) Logger {
	return c.Value(key).(Logger)
}

// ToContext adds a Logger to this context if it supports
// the Setter interface.
func ToContext(c Setter, l Logger) {
	c.Set(key, l)
}

// NewContext instantiate a new context with a kvstore
func NewContext(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, key, l)
}
