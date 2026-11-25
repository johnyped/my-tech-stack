// Package logger is a thin local wrapper around zap, configured with a shared
// log schema (timestamp/level/msg/caller keys, service name
// baked into every line). Pattern copied from ekyc-service's pkg/logger —
// every service keeps a facade like this rather than calling zap directly,
// so logger construction stays a one-line package-level call in main.go.
package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// schemaVersion identifies the log schema version.
const schemaVersion = "1.0"

// Config holds logger initialization parameters.
type Config struct {
	Level       string // "debug", "info", "warn", "error"
	Encoding    string // "json" (prod) or "console" (dev)
	ServiceName string
}

var log *zap.Logger

// Initialize creates and sets the global zap logger. Call once from main.go.
func Initialize(cfg Config) error {
	level, err := parseLevel(cfg.Level)
	if err != nil {
		return fmt.Errorf("invalid log level %q: %w", cfg.Level, err)
	}

	encoding := cfg.Encoding
	if encoding == "" {
		encoding = "json"
	}

	zapCfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Encoding:          encoding,
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		EncoderConfig:     encoderConfig(encoding),
		DisableStacktrace: true,
	}

	l, err := zapCfg.Build()
	if err != nil {
		return fmt.Errorf("build logger: %w", err)
	}

	l = l.With(
		zap.String("schemaVersion", schemaVersion),
		zap.String("serviceName", cfg.ServiceName),
	)

	log = l
	zap.ReplaceGlobals(l)
	return nil
}

// Get returns the global logger instance, falling back to a no-op logger if
// Initialize has not been called (e.g. in a unit test that doesn't need logs).
func Get() *zap.Logger {
	if log == nil {
		return zap.NewNop()
	}
	return log
}

// Sync flushes any buffered log entries. Call via defer in main.go.
func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}

func parseLevel(level string) (zapcore.Level, error) {
	switch level {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info", "":
		return zapcore.InfoLevel, nil
	case "warn":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("unsupported level: %s", level)
	}
}

func encoderConfig(encoding string) zapcore.EncoderConfig {
	cfg := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	if encoding == "console" {
		cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		cfg.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	return cfg
}
