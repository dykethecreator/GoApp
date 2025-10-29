package logger

import "go.uber.org/zap"

// NewLogger creates a new logger.
func NewLogger() (*zap.Logger, error) {
	return zap.NewProduction()
}
