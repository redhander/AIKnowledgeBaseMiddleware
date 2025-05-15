package logger

import (
	"fmt"
	"io"
	"strings"
	"time"
)

type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})

	WithFields(fields Fields) Logger
}

type Fields map[string]interface{}

type logger struct {
	out       io.Writer
	fields    Fields
	callDepth int
}

func New(w io.Writer) Logger {
	return &logger{
		out:       w,
		fields:    make(Fields),
		callDepth: 2,
	}
}

func (l *logger) log(level, message string, args ...interface{}) {
	msg := fmt.Sprintf(message, args...)
	timestamp := time.Now().Format(time.RFC3339)

	var fieldStr strings.Builder
	for k, v := range l.fields {
		fieldStr.WriteString(fmt.Sprintf(" %s=%v", k, v))
	}

	logEntry := fmt.Sprintf("%s [%s] %s%s\n",
		timestamp,
		strings.ToUpper(level),
		msg,
		fieldStr.String(),
	)

	l.out.Write([]byte(logEntry))
}

func (l *logger) Debug(args ...interface{}) {
	l.log("debug", fmt.Sprint(args...))
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.log("debug", format, args...)
}

func (l *logger) Info(args ...interface{}) {
	l.log("info", fmt.Sprint(args...))
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.log("info", format, args...)
}

func (l *logger) Warn(args ...interface{}) {
	l.log("warn", fmt.Sprint(args...))
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.log("warn", format, args...)
}

func (l *logger) Error(args ...interface{}) {
	l.log("error", fmt.Sprint(args...))
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.log("error", format, args...)
}
func (l *logger) Fatalf(format string, args ...interface{}) {
	l.log("fatal", format, args...)
}

func (l *logger) WithFields(fields Fields) Logger {
	newFields := make(Fields)
	for k, v := range l.fields {
		newFields[k] = v
	}
	for k, v := range fields {
		newFields[k] = v
	}

	return &logger{
		out:       l.out,
		fields:    newFields,
		callDepth: l.callDepth,
	}
}
