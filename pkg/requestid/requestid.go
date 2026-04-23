package requestid

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const headerRequestID = "X-Request-ID"

func New(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		rid := c.Request().Header.Get(headerRequestID)
		if rid == "" {
			rid = uuid.NewString()
		}

		c.Request().Header.Set(headerRequestID, rid)
		c.Response().Header().Set(headerRequestID, rid)

		return next(c)
	}
}

func Get(c echo.Context) string {
	return c.Request().Header.Get(headerRequestID)
}
