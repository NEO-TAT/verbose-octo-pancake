package logrux

import (
	"github.com/NEO-TAT/tat_auto_roll_call_service/pkg/env"
	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Logger {
	logger := logrus.New()

	if env.IsDebugMode {
		logger.SetLevel(logrus.DebugLevel)
	}

	logger.SetReportCaller(env.IsDebugMode)

	return logger
}
