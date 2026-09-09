package client

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/buildbeaver/buildbeaver/common/logger"
	"github.com/buildbeaver/buildbeaver/common/models"
	"github.com/buildbeaver/buildbeaver/server/api/rest/documents"
)

// TestCreateArtifact_SendsCorrectContentMD5Header verifies that CreateArtifact computes a real
// Content-MD5 digest of the artifact data and sends it to the server, and that the request body
// received is unaffected by having read it once already to compute the hash (i.e. the reader is
// correctly seeked back to the start before the request is sent).
func TestCreateArtifact_SendsCorrectContentMD5Header(t *testing.T) {
	content := []byte("some artifact content to hash and upload")
	sum := md5.Sum(content)
	expectedHeader := hex.EncodeToString(sum[:])

	var (
		receivedHeader string
		receivedBody   []byte
	)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedHeader = r.Header.Get("Content-MD5")
		var err error
		receivedBody, err = ioutil.ReadAll(r.Body)
		require.NoError(t, err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(&documents.Artifact{})
	}))
	defer server.Close()

	apiClient, err := NewAPIClient([]string{server.URL}, nil, logger.NoOpLogFactory)
	require.NoError(t, err)

	_, err = apiClient.CreateArtifact(
		context.Background(),
		models.JobID{},
		"test-group",
		"test-path.txt",
		bytes.NewReader(content))
	require.NoError(t, err)

	require.Equal(t, expectedHeader, receivedHeader, "expected a real hex Content-MD5 digest of the artifact data, not an empty header")
	require.Equal(t, content, receivedBody, "the reader must be seeked back to the start after hashing so the full body is still sent")
}
