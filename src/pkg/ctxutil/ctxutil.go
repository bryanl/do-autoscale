package ctxutil

import (
	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
)

// LogFromContext extracts a long from a context.Context
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
