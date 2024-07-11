package logger

import (
	"context"
	"os"
	"strings"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

const (
	X_CORRELATION_ID = "X-Correlation-ID"
	CORRELATION_ID   = "Correlation-ID"
)

func Init(level string) {
	// SetFormatter
	log.SetFormatter(&log.JSONFormatter{})

	// SetLevel
	switch strings.ToLower(level) {
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
	log.SetOutput(os.Stdout)
}

func GetLogEntry(ctx context.Context, opName string) *log.Entry {
	var corrID string
	val := ctx.Value(X_CORRELATION_ID)
	if val != nil {
		corrID = val.(string)
	}

	if corrID == "" {
		corrID = uuid.New().String()
	}

	return log.WithFields(log.Fields{
		CORRELATION_ID: corrID,
		"operation":    opName,
	})
}

func GetRequestLogEntry(ctx context.Context, opName string, req interface{}) *log.Entry {
	return GetLogEntry(ctx, opName).WithField("request", req)
}
