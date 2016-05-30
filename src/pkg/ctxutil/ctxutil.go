package ctxutil

import (
	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
)

const (
	KeyEnv = iota
	KeyLog
	KeyDOToken
)

// LogFromContext extracts a log from a context.Context
func LogFromContext(ctx context.Context) *logrus.Entry {
	v := ctx.Value(KeyLog)

	switch v.(type) {
	case *logrus.Entry:
		return v.(*logrus.Entry)
	default:
		logger := logrus.New()
		log := logrus.NewEntry(logger)
		return log
	}
}

// StringFromContext extracts a string from a context.Context
func StringFromContext(ctx context.Context, key string) string {
	s := ctx.Value(key)

	switch s.(type) {
	case string:
		return s.(string)
	default:
		log := LogFromContext(ctx)
		log.WithField("key", key).Warn("context key was not present")
		return ""
	}
}

// IsCurrentEnv checks the context for the current environment name.
func IsCurrentEnv(ctx context.Context, envName string) bool {
	if s, ok := ctx.Value(KeyEnv).(string); ok {
		return s == envName
	}

	return false
}
