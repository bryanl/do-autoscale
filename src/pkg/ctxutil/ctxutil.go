package ctxutil

import (
	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
)

// LogFromContext extracts a long from a context.Context
func LogFromContext(ctx context.Context) *logrus.Entry {
	return ctx.Value("log").(*logrus.Entry)
}
