package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/MaximBayurov/rate-limiter/internal/configuration"
)

type Logger interface {
	Debug(string, ...any)
	Info(string, ...any)
	Warn(string, ...any)
	Error(string, ...any)
}

func New(configs configuration.LoggerConf) (Logger, error) {
	var err error
	if err = os.MkdirAll(filepath.Dir(configs.File), 0o755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	var writer io.Writer
	if configs.File != "" {
		var logFile *os.File
		if logFile, err = os.OpenFile(configs.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644); err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}

		writer = io.MultiWriter(os.Stdout, logFile)
	} else {
		writer = io.MultiWriter(os.Stdout)
	}

	level := &slog.LevelVar{}
	if err = level.UnmarshalText([]byte(configs.Level)); err != nil {
		return nil, fmt.Errorf("failed to parse log level: %w", err)
	}

	var logger Logger = slog.New(slog.NewJSONHandler(writer, &slog.HandlerOptions{
		Level: level,
	}))

	return logger, nil
}
