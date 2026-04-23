package requestlog

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"apigo/pkg/apierror"
	"apigo/pkg/requestid"

	"github.com/labstack/echo/v4"
)

func Completed(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Method == "OPTIONS" {
			return next(c)
		}

		start := time.Now()
		err := next(c)

		status := c.Response().Status
		if err != nil {
			var apiErr *apierror.Error
			var he *echo.HTTPError
			switch {
			case errors.As(err, &apiErr):
				status = apiErr.Status
			case errors.As(err, &he):
				status = he.Code
			default:
				status = 500
			}
		}

		slog.Debug("request completed",
			slog.String("request_id", requestid.Get(c)),
			slog.String("method", c.Request().Method),
			slog.String("uri", c.Request().URL.Path),
			slog.String("client_ip", c.RealIP()),
			slog.String("duration", fmt.Sprintf("%v", time.Since(start))),
			slog.String("host", c.Request().Host),
			slog.String("user_agent", c.Request().UserAgent()),
			slog.Int("status", status),
		)

		return err
	}
}
