package logger

import (
	"io"
	"log"
	"log/slog"
	"main/internal/config"
	"os"
)

// The application sets up one global logger.
// For each incoming request, the middleware then generates a request-specific logger
// and embeds it within the context for that particular chain of operations
//
// more detail: https://dev.to/ilyakaznacheev/where-to-place-logger-in-golang-13o3

const (
	writerStdout = "stdout"
	writerStderr = "stderr"
	writerFile   = "file"
	jsonFormat   = "json"
)

func SetupDefaultLogger(cfg *config.Config) *slog.Logger {
	var writer io.Writer
	switch cfg.Destination {
	case writerStdout:
		writer = os.Stdout
	case writerStderr:
		writer = os.Stderr
	case writerFile:
		writer = NewLogFile(cfg.DestinationFile)
	default:
		writer = nil
		log.Fatalf("Default logger not created. Invalid 'destination' parameter: %s", cfg.Destination)
	}

	logger := NewLogger(
		WithLevel(cfg.Level),
		WithFormat(cfg.Format),
		WithWriter(writer),
		WithAddSource(cfg.AddSource),
		WithSetDefault(true), // Set as default logger for application
	)
	return logger
}

type loggerOptions struct {
	Level      slog.Level
	AddSource  bool
	IsJSON     bool
	SetDefault bool
	Writer     io.Writer
}

type loggerOption func(*loggerOptions)

func NewLogger(opts ...loggerOption) *slog.Logger {
	// Create config by default
	config := &loggerOptions{
		Level:      slog.LevelInfo,
		AddSource:  false,
		IsJSON:     false,
		SetDefault: false,
		Writer:     os.Stdout,
	}
	// Override by custom options
	for _, opt := range opts {
		opt(config)
	}

	// Applying default and custom options
	options := &slog.HandlerOptions{
		AddSource: config.AddSource,
		Level:     config.Level,
	}

	var handler slog.Handler
	if config.IsJSON {
		handler = slog.NewJSONHandler(config.Writer, options)
	} else {
		handler = slog.NewTextHandler(config.Writer, options)
	}

	logger := slog.New(handler)

	if config.SetDefault {
		slog.SetDefault(logger)
	}

	return logger
}

func WithLevel(level string) loggerOption {
	return func(opts *loggerOptions) {
		var sl slog.Level
		err := sl.UnmarshalText([]byte(level))
		if err == nil {
			opts.Level = sl
		} else {
			opts.Level = slog.LevelInfo
		}
	}
}

func WithFormat(format string) loggerOption {
	return func(opts *loggerOptions) {
		opts.IsJSON = format == jsonFormat
	}
}

func WithWriter(writer io.Writer) loggerOption {
	return func(opts *loggerOptions) {
		opts.Writer = writer
	}
}

func WithAddSource(addSource bool) loggerOption {
	return func(opts *loggerOptions) {
		opts.AddSource = addSource
	}
}

func WithSetDefault(setDefault bool) loggerOption {
	return func(opts *loggerOptions) {
		opts.SetDefault = setDefault
	}
}

// Create file for logging
func NewLogFile(path string) *os.File {
	logFile, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("cannot open file: %s", err)
	}
	return logFile
}
