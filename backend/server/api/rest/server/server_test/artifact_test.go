package api_test

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/buildbeaver/buildbeaver/common/models"
	"github.com/buildbeaver/buildbeaver/server/api/rest/client/clienttest"
	"github.com/buildbeaver/buildbeaver/server/app/server_test"
)

// TestArtifactAPI_UploadDownloadRoundTrip exercises a real artifact upload through the REST client
// (not a direct ArtifactService call, which would bypass the client's Content-MD5 computation and
// the server handler's header parsing). It confirms the client now computes and sends a real
// Content-MD5 header rather than the empty string it used to send, that the server accepts it and
// persists the resulting hash against the artifact, and that a downloaded copy matches byte for
// byte.
func TestArtifactAPI_UploadDownloadRoundTrip(t *testing.T) {
	ctx := context.Background()

	app, cleanup, err := server_test.New(server_test.TestConfig(t))
	require.NoError(t, err)
	defer cleanup()
	app.RunnerAPIServer.Start()
	defer app.RunnerAPIServer.Stop(ctx)

	apiClient, clientCert := clienttest.MakeClientCertificateAPIClient(t, app)

	company := server_test.CreateCompanyLegalEntity(t, ctx, app, "", "", "")
	_ = server_test.CreateRunner(t, ctx, app, "", company.ID, clientCert)
	repo := server_test.CreateRepo(t, ctx, app, company.ID)
	commit := server_test.CreateCommit(t, ctx, app, repo.ID, company.ID)

	_, err = app.QueueService.EnqueueBuildFromCommit(ctx, nil, commit, "refs/heads/master", nil)
	require.NoError(t, err)

	job, err := apiClient.Dequeue(ctx)
	require.NoError(t, err)
	require.NotNil(t, job, "expected a job to be available to dequeue")

	content := []byte("some test artifact content, used to verify end-to-end MD5 integrity checking")
	expectedSum := md5.Sum(content)
	expectedHash := hex.EncodeToString(expectedSum[:])

	created, err := apiClient.CreateArtifact(ctx, job.Job.ID, "test-group", "test-artifact.txt", bytes.NewReader(content))
	require.NoError(t, err)
	require.Equal(t, models.HashTypeMD5, created.HashType)
	require.Equal(t, expectedHash, created.Hash, "server-persisted hash should match the actual content, proving the client sent a real Content-MD5 rather than an empty one")
	require.Equal(t, uint64(len(content)), created.Size)

	reader, err := apiClient.GetArtifactData(ctx, created.ID)
	require.NoError(t, err)
	defer reader.Close()
	downloaded, err := ioutil.ReadAll(reader)
	require.NoError(t, err)
	require.Equal(t, content, downloaded)
}
