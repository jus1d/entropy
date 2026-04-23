package ratelimit

import (
	"net/http"

	"apigo/pkg/apierror"

	"github.com/labstack/echo/v4"
)

var errRateLimited = apierror.New(http.StatusTooManyRequests, apierror.CodeRateLimit, "rate limit exceeded", "Slow down and retry after a moment")

func middleware(allow func(ip string) bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !allow(c.RealIP()) {
				return errRateLimited
			}
			return next(c)
		}
	}
}
