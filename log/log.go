package log

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/isomnath/tiny-url/config"

	"github.com/sirupsen/logrus"
)

const (
	redis = "redis"
)

const (
	na          = "na"
	httpRequest = "http.request"
	datastore   = "datastore"
)

var persistenceTypes map[string]string

type Logger struct {
	*logrus.Logger
}

var Log *Logger

type ErrorLogger struct {
	Error error
}

func Setup() {
	level, err := logrus.ParseLevel(config.GetAppLogLevel())
	if err != nil {
		level = logrus.InfoLevel
	}

	persistenceTypes = map[string]string{
		redis: "Redis",
	}

	logrusVars := &logrus.Logger{
		Out:       os.Stderr,
		Hooks:     make(logrus.LevelHooks),
		Formatter: &logrus.JSONFormatter{},
		Level:     level,
	}

	Log = &Logger{logrusVars}
}

func (logger *Logger) getBaseLogEntry() *logrus.Entry {
	return logger.WithFields(
		logrus.Fields{
			"application": map[string]string{
				"name":        config.GetAppName(),
				"version":     config.GetAppVersion(),
				"environment": config.GetAppEnvironment(),
			},
		})
}

func (logger *Logger) baseLogEntry() *logrus.Entry {
	return logger.getBaseLogEntry().WithField("context", na)
}

func (logger *Logger) Fatalf(format string, args ...interface{}) {
	logger.baseLogEntry().
		WithFields(logger.getProcessFields(runtime.Caller(1))).
		Errorf(format, args...)
}

func (logger *Logger) Errorf(format string, args ...interface{}) {
	logger.baseLogEntry().
		WithFields(logger.getProcessFields(runtime.Caller(1))).
		Errorf(format, args...)
}

func (logger *Logger) Infof(format string, args ...interface{}) {
	logger.baseLogEntry().
		WithFields(logger.getProcessFields(runtime.Caller(1))).
		Infof(format, args...)
}

func (logger *Logger) Warnf(format string, args ...interface{}) {
	logger.baseLogEntry().
		WithFields(logger.getProcessFields(runtime.Caller(1))).
		Warnf(format, args...)
}

func (logger *Logger) httpRequestLogEntry(r *http.Request) *logrus.Entry {
	return logger.getBaseLogEntry().
		WithFields(logrus.Fields{
			"context": httpRequest,
			"request": map[string]interface{}{
				"path":           r.URL.Path,
				"method":         r.Method,
				"host":           r.Host,
				"remote_address": r.RemoteAddr,
			},
		})
}

func (logger *Logger) HTTPStatInfo(r *http.Request, startTime, responseTime time.Time, statusCode int) {
	latency := responseTime.Sub(startTime)
	fields := logrus.Fields{
		"context": httpRequest,
		"request": map[string]interface{}{
			"host":           r.Host,
			"method":         r.Method,
			"path":           r.URL.Path,
			"start_time":     startTime.Format(time.RFC3339),
			"remote_address": r.RemoteAddr,
		},
		"response": map[string]interface{}{
			"end_time": responseTime.Format(time.RFC3339),
			"latency":  fmt.Sprintf("%d ms", latency.Milliseconds()),
			"status":   statusCode,
		},
		"forwarded_headers": r.Header.Get("X_FORWARDED-FOR"),
	}
	logger.httpRequestLogEntry(r).
		WithFields(fields).
		WithFields(logger.getProcessFields(runtime.Caller(1))).
		Infof("http log")
}

func (logger *Logger) HTTPErrorf(r *http.Request, format string, args ...interface{}) {
	logger.httpRequestLogEntry(r).
		WithFields(logger.getProcessFields(runtime.Caller(1))).
		Errorf(format, args...)
}

func (logger *Logger) HTTPInfof(r *http.Request, format string, args ...interface{}) {
	logger.httpRequestLogEntry(r).
		WithFields(logger.getProcessFields(runtime.Caller(1))).
		Infof(format, args...)
}

func (logger *Logger) HTTPWarnf(r *http.Request, format string, args ...interface{}) {
	logger.httpRequestLogEntry(r).
		WithFields(logger.getProcessFields(runtime.Caller(1))).
		Warnf(format, args...)
}

func (logger *Logger) persistenceLogEntry(dbType string) *logrus.Entry {
	persistenceType := persistenceTypes[dbType]
	return logger.getBaseLogEntry().
		WithFields(logrus.Fields{
			"context": datastore,
			"type":    persistenceType,
		})
}

func (logger *Logger) RedisErrorf(format string, args ...interface{}) {
	logger.persistenceLogEntry(redis).
		WithFields(logger.getProcessFields(runtime.Caller(1))).
		Errorf(format, args...)
}

func (logger *Logger) RedisInfof(format string, args ...interface{}) {
	logger.persistenceLogEntry(redis).
		WithFields(logger.getProcessFields(runtime.Caller(1))).
		Infof(format, args...)
}

func (logger *Logger) RedisWarnf(format string, args ...interface{}) {
	logger.persistenceLogEntry(redis).
		WithFields(logger.getProcessFields(runtime.Caller(1))).
		Warnf(format, args...)
}

func (logger *Logger) getProcessFields(pc uintptr, file string, line int, ok bool) logrus.Fields {
	var fileName, fn string
	if !ok {
		fileName = "unknown"
		fn = "unknown"
	} else {
		fileName = file[strings.LastIndex(file, "/")+1:] + ":" + strconv.Itoa(line)
		fnName := runtime.FuncForPC(pc).Name()
		fn = fnName[strings.LastIndex(fnName, ".")+1:]
	}

	return logrus.Fields{
		"file":     fileName,
		"function": fn,
	}
}
