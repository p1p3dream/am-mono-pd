package gconf

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"abodemine/lib/errors"
	"abodemine/lib/val"
)

// Sentry represents the sentry configuration
type Sentry struct {
	DSN              string  `json:"dsn,omitempty" yaml:"dsn,omitempty"`
	TracesSampleRate float64 `json:"traces_sample_rate,omitempty" yaml:"traces_sample_rate,omitempty"`
	EnableTracing    bool    `json:"enable_tracing,omitempty" yaml:"enable_tracing,omitempty"`
}

func LoadZerolog(logLevel string, noLogColor bool) error {
	zerologConfig, err := GetZerologConfig(logLevel, noLogColor)
	if err != nil {
		return errors.Forward(err, "e9420c7c-6e89-44d2-b823-e2e339bdd3b1")
	}

	zerolog.SetGlobalLevel(zerologConfig.Level)

	var writer io.Writer

	if zerologConfig.ConsoleWriter.NoColor {
		writer = zerologConfig.ConsoleWriter.Out
	} else {
		writer = zerologConfig.ConsoleWriter
	}

	log.Logger = zerolog.New(writer).
		With().
		Timestamp().
		Logger()

	return nil
}

type Zerolog struct {
	ConsoleWriter zerolog.ConsoleWriter
	Level         zerolog.Level
}

func GetZerologConfig(logLevel string, noLogColor bool) (*Zerolog, error) {
	logLevel = val.Coalesce(
		os.Getenv("ZEROLOG_LEVEL"),
		strings.ToLower(strings.TrimSpace(logLevel)),
		"info",
	)

	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		return nil, &errors.Object{
			Id:     "adc89bdd-7e41-4ff9-b7ec-962422a9cd1d",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Invalid zerolog level.",
			Cause:  err.Error(),
		}
	}

	consoleWriter := zerolog.ConsoleWriter{
		TimeFormat: time.RFC3339Nano,
		Out:        os.Stderr,
	}

	if noLogColor {
		consoleWriter.NoColor = true
	}

	return &Zerolog{
		ConsoleWriter: consoleWriter,
		Level:         level,
	}, nil
}
