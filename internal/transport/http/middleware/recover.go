package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime"

	"apigo/pkg/requestid"

	"github.com/labstack/echo/v4"
)

func Recover(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, 2048)
				n := runtime.Stack(buf, false)

				slog.Error("panic recovered",
					slog.String("request_id", requestid.Get(c)),
					slog.String("error", fmt.Sprintf("%v", r)),
					slog.String("stack", string(buf[:n])),
				)

				_ = c.JSON(http.StatusInternalServerError, map[string]string{
					"message":    "internal server error",
					"request_id": requestid.Get(c),
				})
			}
		}()

		return next(c)
	}
}
