package clog

import (
	"context"
	validation "github.com/itgelo/ozzo-validation/v4"
	"github.com/minhho2511/elotusteam-test/internal/middleware"
	"time"
)

const (
	// InfoLevel is the default logging priority.
	InfoLevel = "info"

	// WarnLevel logs are more important than Info, but don't need individual human review.
	WarnLevel = "warn"

	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel = "error"

	//DebugLevel Anything else, i.e. too verbose to be included in INFO level.
	DebugLevel = "debug"
)

type Writer interface {
	Printf(logCollection *LogCollection)
	Log(keyvals ...interface{}) error
}

type LogCollection struct {
	DateTime   time.Time
	Level      string
	TraceID    string
	Message    string
	Caller     string
	StackTrace string
}

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(err error)
	Debug(err error)
	Handle(ctx context.Context, err error)
	WithContext(ctx context.Context) Logger
	Log(keyvals ...interface{}) error
}

type logger struct {
	Writer  Writer
	Context context.Context
}

func NewLogger(writer Writer) Logger {
	return &logger{
		Writer: writer,
	}
}

func (l *logger) WithContext(ctx context.Context) Logger {
	lg := *l
	lg.Context = ctx
	return &lg
}

func (l *logger) Info(msg string) {
	l.Writer.Printf(l.format(msg, InfoLevel))
}

func (l *logger) Warn(msg string) {
	l.Writer.Printf(l.format(msg, WarnLevel))
}

func (l *logger) Error(err error) {
	l.Writer.Printf(l.format(err.Error(), ErrorLevel))
}

func (l *logger) Debug(err error) {
	l.Writer.Printf(l.format(err.Error(), DebugLevel))
}

func (l *logger) Log(keyvals ...interface{}) error {
	traceId := l.traceId()
	if len(traceId) > 0 {
		keyvals = append([]interface{}{"trace-id", traceId}, keyvals...)
	}
	return l.Writer.Log(keyvals...)
}

func (l *logger) format(msg string, level string) *LogCollection {
	return &LogCollection{
		DateTime: time.Now(),
		Level:    level,
		TraceID:  l.traceId(),
		Message:  msg,
	}
}

func (l *logger) traceId() string {
	var traceId string
	if l.Context != nil {
		if ctxTraceID, ok := l.Context.Value(middleware.TraceIDContextKey).(string); ok {
			traceId = ctxTraceID
		}
	}
	return traceId
}

// Handle - implement for ServerErrorHandler
func (l *logger) Handle(_ context.Context, err error) {
	switch err.(type) {
	default:
		l.Error(err)
	case validation.Errors:
		l.Info(err.Error())
	}
}
