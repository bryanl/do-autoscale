package ctxutil

import (
	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
)

// LogFromContext extracts a log from a context.Context
func LogFromContext(ctx context.Context) *logrus.Entry {
	v := ctx.Value("log")

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
