package logger

import (
	"go.uber.org/zap"
)

var log *zap.SugaredLogger

func init() {
	// Development logger (pretty print)
	zapLogger, _ := zap.NewDevelopment()
	log = zapLogger.Sugar()
}

// Info logs an info level message
func Info(msg string, args ...interface{}) {
	log.Infof(msg, args...)
}

// Error logs an error level message
func Error(msg string, args ...interface{}) {
	log.Errorf(msg, args...)
}

// Example: structured logging with fields
func InfoWithFields(msg string, fields map[string]interface{}) {
	log.Infow(msg, fieldsToZap(fields)...)
}

func ErrorWithFields(msg string, fields map[string]interface{}) {
	log.Errorw(msg, fieldsToZap(fields)...)
}

// helper to convert map to zap.Fields
func fieldsToZap(fields map[string]interface{}) []interface{} {
	zapFields := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		zapFields = append(zapFields, k, v)
	}
	return zapFields
}
