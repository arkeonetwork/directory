// Structured logging with field annotation support.
// Details about where the log was called from are included.

package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Fields is an alias to logrus.Fields
type Fields = logrus.Fields

// Logger wraps a logrus.FieldLogger to provide all standard logging functionality
type Logger interface {
	logrus.FieldLogger
}

const timestampFormat = "2006-01-02T15:04:05.999Z07:00"

var logger Logger

func init() {
	logger = logrus.StandardLogger()

	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{TimestampFormat: timestampFormat})
	if os.Getenv("ARKEO_DIR_JSON_LOGS") == "true" {
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: timestampFormat,
			PrettyPrint:     false,
		})
	}
}

// WithFields adds field annotations to the logger instance
func WithFields(fields Fields) Logger {
	return logger.WithFields(fields)
}

// WithoutFields uses the default logger with no extra field annotations
func WithoutFields() Logger {
	return logger
}
