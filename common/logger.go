package common

import "go.uber.org/zap"

var logger *zap.Logger

// ConfigureLogger configure logger instance.
func ConfigureLogger() error {
	var err error
	// TODO Configure logger.
	logger, err = zap.NewDevelopment()
	if err != nil {
		return err
	}

	return nil
}

// GetLogger return a logger instance.
func GetLogger() *zap.Logger {
	return logger
}
