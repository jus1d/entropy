package http

import (
	"log/slog"
	"net/http"

	"apigo/pkg/apierror"
	"apigo/pkg/apiresponse"
	"apigo/pkg/log/sl"
	"apigo/pkg/requestid"

	"github.com/labstack/echo/v4"
)

func HTTPErrorHandler(err error, c echo.Context) {
	reqID := requestid.Get(c)

	if he, ok := err.(*echo.HTTPError); ok {
		if he.Code == http.StatusNotFound {
			_ = apiresponse.Error(c, http.StatusNotFound, apierror.CodeNotFound, "resource not found", "Check the URL and HTTP method")
			return
		}
		if he.Code == http.StatusMethodNotAllowed {
			_ = apiresponse.Error(c, http.StatusMethodNotAllowed, apierror.CodeMethodNotAllowed, "method not allowed", "Check the HTTP method and try again")
			return
		}
	}

	if ae, ok := err.(*apierror.Error); ok {
		slog.Debug("responded with API error",
			sl.Err(err),
			slog.String("request_id", reqID),
			slog.Int("status", ae.Status),
			slog.String("error_code", string(ae.Code)),
			slog.String("hint", ae.Hint),
			slog.String("method", c.Request().Method),
			slog.String("uri", c.Request().URL.Path),
		)
		_ = apiresponse.Error(c, ae.Status, ae.Code, ae.Message, ae.Hint)
		return
	}

	slog.Error("unhandled error",
		sl.Err(err),
		slog.String("request_id", reqID),
		slog.String("method", c.Request().Method),
		slog.String("uri", c.Request().URL.Path),
		slog.String("client_ip", c.RealIP()),
	)
	_ = apiresponse.Error(c, http.StatusInternalServerError, apierror.CodeInternal, "internal server error", "Try again later or contact support")
}
