package gophoenixsocketclient

import (
	"context"
	log "github.com/sirupsen/logrus"
	"os"
	"runtime"
)

// Logger logging system interface
type Logger interface {
	Info(service string, method string, args ...interface{})
	Warn(service string, method string, args ...interface{})
	Debug(service string, method string, args ...interface{})
	Error(service string, method string, args ...interface{})
	Printf(format string, args ...interface{})
	Level(string) Logger
	Format(string) Logger
}

type logger struct {
	l       *log.Logger
	e       *log.Entry
	ctx     context.Context
	version string
}

// NewLogger initialize logging system
func NewLogger(ctx context.Context, version string) Logger {
	l := new(logger)
	l.ctx = ctx
	l.version = version
	l.init()

	return l
}

func (l *logger) init() {
	l1 := log.New()
	l1.SetFormatter(&log.TextFormatter{})
	l1.SetOutput(os.Stdout)

	l.l = l1
	l.e = l1.WithContext(l.ctx).WithFields(log.Fields{
		"name":    "transferchain-desktop",
		"version": l.version,
	})
}

// Format ..
func (l *logger) Format(format string) Logger {
	switch format {
	case "json":
		l.l.SetFormatter(&log.JSONFormatter{})
		break
	case "text":
		l.l.SetFormatter(&log.TextFormatter{})
		break
	}

	return l
}

// Level ..
func (l *logger) Level(level string) Logger {
	lNew := new(logger)
	lNew.init()
	switch level {
	case "info":
		l.l.SetLevel(log.InfoLevel)
		break
	case "debug":
		l.l.SetLevel(log.DebugLevel)
		break
	case "warning":
		l.l.SetLevel(log.WarnLevel)
		break
	case "error":
		l.l.SetLevel(log.ErrorLevel)
		break
	}
	return lNew
}

// Printf ..
func (l *logger) Printf(format string, args ...interface{}) {
	_, _, no, _ := runtime.Caller(1)
	l.e.WithFields(log.Fields{
		"line": no,
	}).Printf(format, args...)
}

// Info ..
func (l logger) Info(service string, method string, args ...interface{}) {
	_, _, no, _ := runtime.Caller(1)
	l.e.WithFields(log.Fields{
		"service": service,
		"method":  method,
		"line":    no,
	}).Info(args...)
}

// Warn ..
func (l logger) Warn(service string, method string, args ...interface{}) {
	_, _, no, _ := runtime.Caller(1)
	l.e.WithFields(log.Fields{
		"service": service,
		"method":  method,
		"line":    no,
	}).Warn(args...)
}

// Debug ..
func (l logger) Debug(service string, method string, args ...interface{}) {
	_, _, no, _ := runtime.Caller(1)
	l.e.WithFields(log.Fields{
		"service": service,
		"method":  method,
		"line":    no,
	}).Debug(args...)
}

// Error ..
func (l logger) Error(service string, method string, args ...interface{}) {
	_, _, no, _ := runtime.Caller(1)
	l.e.WithFields(log.Fields{
		"service": service,
		"method":  method,
		"line":    no,
	}).Error(args...)
}
