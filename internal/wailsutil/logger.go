package wailsutil

import "github.com/iver-wharf/wharf-core/v2/pkg/logger"

type Logger struct {
	WharfLogger logger.Logger
}

func (log Logger) Print(message string) {
	log.WharfLogger.Info().Message(message)
}

func (log Logger) Trace(message string) {
	log.WharfLogger.Debug().Message(message)
}

func (log Logger) Debug(message string) {
	log.WharfLogger.Debug().Message(message)
}

func (log Logger) Info(message string) {
	log.WharfLogger.Info().Message(message)
}

func (log Logger) Warning(message string) {
	log.WharfLogger.Warn().Message(message)
}

func (log Logger) Error(message string) {
	log.WharfLogger.Error().Message(message)
}

func (log Logger) Fatal(message string) {
	log.WharfLogger.Error().Message(message)
}
