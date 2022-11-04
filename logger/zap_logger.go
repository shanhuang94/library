package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/opentracing/opentracing-go"
	opentracinglog "github.com/opentracing/opentracing-go/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	jae "neuroxess-cloud/library/jaeger"
	"neuroxess-cloud/library/logger/utils"
	"os"
	"path/filepath"
	"sync"
)

var ZapLogger *zap.Logger
var once sync.Once

const (
	Level = "level"
)

type LogLevel string

var panic LogLevel = "panic"
var fatal LogLevel = "fatal"
var error LogLevel = "error"
var debug LogLevel = "debug"
var warn LogLevel = "warn"
var info LogLevel = "info"

func InitLogger(path string, level string, isDebug bool) {
	log.Printf("initLogger path %s, level %s, isDebug %t\n", path, level, isDebug)

	logBasePath, err := filepath.Abs(filepath.Dir(path))
	if _, err := os.Stat(logBasePath); err != nil {
		err = os.MkdirAll(logBasePath, 0711)
	}

	var js string
	if isDebug {
		js = fmt.Sprintf(`{
      "level": "%s",
      "encoding": "json",
      "outputPaths": ["stdout"],
      "errorOutputPaths": ["stdout"]
      }`, level)
	} else {
		js = fmt.Sprintf(`{
      "level": "%s",
      "encoding": "json",
      "outputPaths": ["%s"],
      "errorOutputPaths": ["%s"]
      }`, level, path, path)
	}

	var cfg zap.Config
	if err := json.Unmarshal([]byte(js), &cfg); err != nil {
		log.Panicln(err)
	}

	cfg.EncoderConfig = zap.NewProductionEncoderConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.DisableStacktrace = true
	cfg.DisableCaller = true

	ZapLogger, err = cfg.Build()
	if err != nil {
		log.Println("init logger error: ", err)
	}
}

func GetLogger() *zap.Logger {
	return ZapLogger
}

// panic
// Panic logs on panic level and trace based on the context span if it exists
func Panicc(ctx context.Context, log string, fields ...zapcore.Field) {
	span := opentracing.SpanFromContext(ctx)
	logSpan(span, panic, log, fields...)
	fields = addFields(fields, span, panic)
	GetLogger().Panic(log, fields...)
}

// fatal
// Fatalc logs on fatal level and trace based on the context span if it exists
func Fatalc(ctx context.Context, log string, fields ...zapcore.Field) {
	span := opentracing.SpanFromContext(ctx)
	logSpan(span, fatal, log, fields...)
	fields = addFields(fields, span, fatal)
	GetLogger().Fatal(log, fields...)
}

// error
// Errorc logs on error level and trace based on the context span if it exists
func Errorc(ctx context.Context, log string, fields ...zapcore.Field) {
	span := opentracing.SpanFromContext(ctx)
	logSpan(span, error, log, fields...)
	fields = addFields(fields, span, error)
	GetLogger().Error(log, fields...)
}

// warn
// Warnc logs on info level and trace based on the context span if it exists
func Warnc(ctx context.Context, log string, fields ...zapcore.Field) {
	span := opentracing.SpanFromContext(ctx)
	logSpan(span, warn, log, fields...)
	fields = addFields(fields, span, warn)
	GetLogger().Warn(log, fields...)
}

// info
// Infoc logs on info level and trace based on the context span if it exists
func Infoc(ctx context.Context, log string, fields ...zapcore.Field) {
	span := opentracing.SpanFromContext(ctx)
	logSpan(span, info, log, fields...)
	//add traceid spanid
	fields = addFields(fields, span, info)
	GetLogger().Info(log, fields...)
}

// debug
// Debugc logs on debug level and trace based on the context span if it exists
func Debugc(ctx context.Context, log string, fields ...zapcore.Field) {
	span := opentracing.SpanFromContext(ctx)
	logSpan(span, debug, log, fields...)
	fields = addFields(fields, span, debug)
	GetLogger().Debug(log, fields...)
}

// fields 增加日志级别 tarceid spanid
func addFields(fields []zapcore.Field, span opentracing.Span, level LogLevel) []zapcore.Field {
	traceid, spanid := jae.GetTraceIdAndSpanId(span)
	fields = append(fields, zap.String("traceid", traceid), zap.String("spanid", spanid), zap.String("level", string(level)))
	return fields
}

// span增加tag event logfields
func logSpan(span opentracing.Span, level LogLevel, log string, fields ...zapcore.Field) {
	if span == nil {
		GetLogger().Info("no span " + log)
		return
	}
	if level == error {
		span.SetTag(Level, string(level))
	}
	traceField := make([]opentracinglog.Field, len(fields)+2)
	traceField[0] = opentracinglog.Event(log)
	traceField[1] = opentracinglog.String("level", string(level))
	for i, v := range fields {
		traceField[i+2] = utils.ZapFieldToOpentracing(v)
	}
	span.LogFields(traceField...)
}

// panic
// Panic logs on panic level
func Panic(log string, fields ...zapcore.Field) {
	GetLogger().Panic(log, fields...)
}

// fatal
// Fatal logs on fatal level
func Fatal(log string, fields ...zapcore.Field) {
	GetLogger().Fatal(log, fields...)
}

// error
// Error logs on error level
func Error(log string, fields ...zapcore.Field) {
	GetLogger().Error(log, fields...)
}

// warn
// Warn logs on info level
func Warn(log string, fields ...zapcore.Field) {
	GetLogger().Warn(log, fields...)
}

// info
// Info logs on info level
func Info(log string, fields ...zapcore.Field) {
	GetLogger().Info(log, fields...)
}

// debug
// Debug logs on debug level
func Debug(log string, fields ...zapcore.Field) {
	GetLogger().Debug(log, fields...)
}
