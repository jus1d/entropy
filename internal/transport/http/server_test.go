package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"entropy/internal/config"
	apphttp "entropy/internal/transport/http"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestServer() *httptest.Server {
	cfg := &config.Config{
		Env: config.EnvLocal,
		Server: config.Server{
			Address:     ":0",
			Timeout:     4 * time.Second,
			IdleTimeout: 60 * time.Second,
		},
	}
	srv := apphttp.NewServer(cfg, nil)
	return httptest.NewServer(srv.Handler())
}

func TestLiveness(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/liveness")
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestReadiness(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/readiness")
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestNotFound(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/nonexistent")
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestRequestID(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/liveness")
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck

	assert.NotEmpty(t, resp.Header.Get("X-Request-ID"))
}
