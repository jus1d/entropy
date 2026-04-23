package v1

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"entropy/internal/transport/http/v1/response"
	"entropy/pkg/apierror"
	"entropy/pkg/apiresponse"
	"entropy/pkg/log/sl"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (r *Router) getLogset(c echo.Context) error {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return apiresponse.Error(c, http.StatusBadRequest, apierror.CodeInvalidRequest, "invalid uuid", "Provide a valid UUID")
	}

	logs, err := r.logsetService.Get(c.Request().Context(), id)
	if err != nil {
		slog.Error("failed to get logset", sl.Err(err))
		return apiresponse.Error(c, http.StatusInternalServerError, apierror.CodeInternal, "failed to get logset", "")
	}

	return apiresponse.Resource(c, http.StatusOK, response.GetLogset{
		UUID:  id,
		Count: len(logs),
		Logs:  logs,
	})
}

func (r *Router) createLogset(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		slog.Debug("missing file in request", sl.Err(err))
		return apiresponse.Error(c, http.StatusBadRequest, apierror.CodeInvalidRequest, "missing file", "Upload a file using the 'file' form field")
	}

	src, err := file.Open()
	if err != nil {
		slog.Error("failed to open uploaded file", sl.Err(err))
		return apiresponse.Error(c, http.StatusInternalServerError, apierror.CodeInternal, "failed to read file", "")
	}
	defer src.Close() //nolint:errcheck

	data, err := io.ReadAll(src)
	if err != nil {
		slog.Error("failed to read uploaded file", sl.Err(err))
		return apiresponse.Error(c, http.StatusInternalServerError, apierror.CodeInternal, "failed to read file", "")
	}

	var logs []map[string]any
	dec := json.NewDecoder(bytes.NewReader(data))
	for {
		var obj map[string]any
		if err := dec.Decode(&obj); err != nil {
			if err == io.EOF {
				break
			}
			slog.Debug("invalid JSON in uploaded file", sl.Err(err))
			return apiresponse.Error(c, http.StatusUnprocessableEntity, apierror.CodeInvalidRequest, "file contains invalid JSON", "Ensure the file contains a sequence of valid JSON objects")
		}
		logs = append(logs, obj)
	}

	if len(logs) == 0 {
		return apiresponse.Error(c, http.StatusUnprocessableEntity, apierror.CodeInvalidRequest, "file contains no JSON objects", "Ensure the file contains at least one JSON object")
	}

	logsetID, err := r.logsetService.Ingest(c.Request().Context(), logs)
	if err != nil {
		slog.Error("failed to ingest logs", sl.Err(err))
		return apiresponse.Error(c, http.StatusInternalServerError, apierror.CodeInternal, "failed to ingest logs", "")
	}

	return apiresponse.Resource(c, http.StatusOK, response.CreateLogset{
		LogsetID: logsetID,
		Count:    len(logs),
	})
}
