package v1_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	v1 "entropy/internal/transport/http/v1"
	"entropy/internal/transport/http/v1/mock"
	"entropy/pkg/apiresponse"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func setupRouter(t *testing.T) (*mock.MockLogsetService, *echo.Echo) {
	t.Helper()
	ctrl := gomock.NewController(t)
	svc := mock.NewMockLogsetService(ctrl)

	e := echo.New()
	g := e.Group("/api/v1")
	r := v1.NewRouter(svc)
	r.Register(g)

	return svc, e
}

func createFileRequest(t *testing.T, content string) *http.Request {
	t.Helper()
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("file", "logs.json")
	require.NoError(t, err)
	_, err = io.WriteString(part, content)
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	req := httptest.NewRequest(http.MethodPost, "/api/v1/logsets", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func TestCreateLogset_Success(t *testing.T) {
	svc, e := setupRouter(t)

	logsetID := uuid.New()
	svc.EXPECT().
		Ingest(gomock.Any(), gomock.Any()).
		Return(logsetID, nil)

	req := createFileRequest(t, `{"level":"info","msg":"hello"}`)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp apiresponse.Res
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "resource", resp.Kind)

	data, ok := resp.Data.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, logsetID.String(), data["logset_id"])
	assert.Equal(t, float64(1), data["count"])
}

func TestCreateLogset_MultipleLogs(t *testing.T) {
	svc, e := setupRouter(t)

	logsetID := uuid.New()
	svc.EXPECT().
		Ingest(gomock.Any(), gomock.Len(3)).
		Return(logsetID, nil)

	content := `{"a":1}
{"b":2}
{"c":3}`
	req := createFileRequest(t, content)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp apiresponse.Res
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	data := resp.Data.(map[string]any)
	assert.Equal(t, float64(3), data["count"])
}

func TestCreateLogset_MissingFile(t *testing.T) {
	_, e := setupRouter(t)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/logsets", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp apiresponse.Err
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "error", resp.Kind)
	assert.Equal(t, "missing file", resp.Error.Message)
}

func TestCreateLogset_InvalidJSON(t *testing.T) {
	_, e := setupRouter(t)

	req := createFileRequest(t, `not json at all`)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	var resp apiresponse.Err
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "file contains invalid JSON", resp.Error.Message)
}

func TestCreateLogset_EmptyFile(t *testing.T) {
	_, e := setupRouter(t)

	req := createFileRequest(t, ``)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	var resp apiresponse.Err
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "file contains no JSON objects", resp.Error.Message)
}

func TestGetLogset_Success(t *testing.T) {
	svc, e := setupRouter(t)

	id := uuid.New()
	logs := []map[string]any{
		{"level": "info", "msg": "hello"},
		{"level": "error", "msg": "world"},
	}
	svc.EXPECT().
		Get(gomock.Any(), id).
		Return(logs, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/logsets/"+id.String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp apiresponse.Res
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "resource", resp.Kind)

	data, ok := resp.Data.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, id.String(), data["uuid"])
	assert.Equal(t, float64(2), data["count"])

	respLogs, ok := data["logs"].([]any)
	require.True(t, ok)
	assert.Len(t, respLogs, 2)
}

func TestGetLogset_EmptyLogs(t *testing.T) {
	svc, e := setupRouter(t)

	id := uuid.New()
	svc.EXPECT().
		Get(gomock.Any(), id).
		Return(nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/logsets/"+id.String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp apiresponse.Res
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	data := resp.Data.(map[string]any)
	assert.Equal(t, float64(0), data["count"])
}

func TestGetLogset_InvalidUUID(t *testing.T) {
	_, e := setupRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/logsets/not-a-uuid", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp apiresponse.Err
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "invalid uuid", resp.Error.Message)
}

func TestGetLogset_ServiceError(t *testing.T) {
	svc, e := setupRouter(t)

	id := uuid.New()
	svc.EXPECT().
		Get(gomock.Any(), id).
		Return(nil, fmt.Errorf("clickhouse down"))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/logsets/"+id.String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp apiresponse.Err
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "failed to get logset", resp.Error.Message)
}

func TestCreateLogset_IngestError(t *testing.T) {
	svc, e := setupRouter(t)

	svc.EXPECT().
		Ingest(gomock.Any(), gomock.Any()).
		Return(uuid.UUID{}, fmt.Errorf("clickhouse down"))

	req := createFileRequest(t, `{"level":"error"}`)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp apiresponse.Err
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "failed to ingest logs", resp.Error.Message)
}
