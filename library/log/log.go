package log

import (
	"fmt"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

func MustBuildLogR() logr.Logger {
	var log logr.Logger

	zapConfig := zap.NewDevelopmentConfig()
	zapLog, err := zapConfig.Build()
	if err != nil {
		panic(fmt.Sprintf("who watches the watchmen (%v)?", err))
	}
	log = zapr.NewLogger(zapLog)

	return log
}
