package config

import (
	"github.com/hashicorp/go-hclog"
)

// NewLogger returns a new logger instance
func NewLogger() hclog.Logger {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "procedure-service",
		Level: hclog.LevelFromString("DEBUG"),
	})

	return logger
}
