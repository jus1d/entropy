package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func liveness(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}

func readiness(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}
