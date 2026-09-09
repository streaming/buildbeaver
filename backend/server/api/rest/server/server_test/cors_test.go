package api_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/buildbeaver/buildbeaver/server/app/server_test"
)

// requestWithOrigin issues a GET request to the Core API's public root document endpoint, as if
// from a browser page on the given origin, and returns the Access-Control-Allow-Origin response
// header (empty if not present).
func requestWithOrigin(t *testing.T, serverURL string, origin string) string {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, serverURL+"/api/v1/", nil)
	require.NoError(t, err)
	req.Header.Set("Origin", origin)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	return resp.Header.Get("Access-Control-Allow-Origin")
}

// TestCORS_DisabledByDefault confirms that the App/Core API does not send permissive CORS headers
// unless explicitly configured to, closing a gap where CORS was previously always enabled
// (allowing localhost/127.0.0.1 origins with credentials) regardless of build/deployment mode.
func TestCORS_DisabledByDefault(t *testing.T) {
	ctx := context.Background()
	config := server_test.TestConfig(t)
	require.False(t, config.CoreAPIConfig.EnableDevCORS, "CORS should be disabled by default")

	app, cleanup, err := server_test.New(config)
	require.NoError(t, err)
	defer cleanup()
	app.CoreAPIServer.Start()
	defer app.CoreAPIServer.Stop(ctx)

	allowOrigin := requestWithOrigin(t, app.CoreAPIServer.GetServerURL(), "http://localhost:3000")
	require.Empty(t, allowOrigin, "no Access-Control-Allow-Origin header should be sent when CORS is not enabled")
}

// TestCORS_EnabledWhenConfigured confirms that setting EnableDevCORS does turn on the permissive
// CORS handling intended for local frontend development.
func TestCORS_EnabledWhenConfigured(t *testing.T) {
	ctx := context.Background()
	config := server_test.TestConfig(t)
	config.CoreAPIConfig.EnableDevCORS = true

	app, cleanup, err := server_test.New(config)
	require.NoError(t, err)
	defer cleanup()
	app.CoreAPIServer.Start()
	defer app.CoreAPIServer.Stop(ctx)

	allowOrigin := requestWithOrigin(t, app.CoreAPIServer.GetServerURL(), "http://localhost:3000")
	require.Equal(t, "http://localhost:3000", allowOrigin)
}
