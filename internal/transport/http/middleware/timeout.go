package middleware

import (
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	"time"
)

func Timeout(duration time.Duration) echo.MiddlewareFunc {
	return echomw.TimeoutWithConfig(echomw.TimeoutConfig{
		Timeout: duration,
	})
}
