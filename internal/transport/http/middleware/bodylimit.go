package middleware

import (
	echomw "github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo/v4"
)

// BodyLimit returns a middleware that limits request body size.
// limit uses Echo's size format: "1M", "512K", "2G", etc.
func BodyLimit(limit string) echo.MiddlewareFunc {
	return echomw.BodyLimit(limit)
}
