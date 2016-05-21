package echologger

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

// New returns a new middleware handler with a default name and logger
func New() echo.MiddlewareFunc {
	return NewWithName("web")
}

// NewWithName returns a new middleware handler with the specified name
func NewWithName(name string) echo.MiddlewareFunc {
	e := logrus.NewEntry(logrus.StandardLogger())
	return NewWithNameAndLogger(name, e)
}

// NewWithNameAndLogger returns a new middleware handler with the specified name
// and logger
func NewWithNameAndLogger(name string, l *logrus.Entry) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			entry := l.WithFields(logrus.Fields{
				"request": c.Request().URI(),
				"method":  c.Request().Method(),
				"remote":  c.Request().RemoteAddress(),
			})

			if reqID := c.Request().Header().Get("X-Request-Id"); reqID != "" {
				entry = entry.WithField("request_id", reqID)
			}

			entry.Info("started handling request")

			if err := next(c); err != nil {
				c.Error(err)
			}

			latency := time.Since(start)

			entry.WithFields(logrus.Fields{
				"status":      c.Response().Status(),
				"text_status": http.StatusText(c.Response().Status()),
				"took":        latency,
				fmt.Sprintf("measure#%s.latency", name): latency.Nanoseconds(),
			}).Info("completed handling request")

			return nil
		}
	}
}
