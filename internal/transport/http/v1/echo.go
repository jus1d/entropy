package v1

import (
	"apigo/pkg/apierror"
	"apigo/pkg/apiresponse"
	"apigo/pkg/log/sl"
	"apigo/pkg/validate"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

func echoHandler(c echo.Context) error {
	var body map[string]any
	if err := validate.Bind(c, &body); err != nil {
		slog.Debug("invalid request body", sl.Err(err))
		return apiresponse.Error(c, http.StatusBadRequest, apierror.CodeInvalidRequest, "invalid request body", "Ensure the request body is valid JSON")
	}

	return apiresponse.Resource(c, http.StatusOK, body)
}
