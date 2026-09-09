package runner

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/buildbeaver/buildbeaver/common/models"
	"github.com/buildbeaver/buildbeaver/runner/logging"
	"github.com/buildbeaver/buildbeaver/server/api/rest/documents"
)

// fakeArtifactDataAPIClient is a minimal APIClient stub that only implements GetArtifactData,
// for testing ArtifactManager.downloadArtifact in isolation. Every other method panics if called.
type fakeArtifactDataAPIClient struct {
	APIClient
	data []byte
}

func (f *fakeArtifactDataAPIClient) GetArtifactData(ctx context.Context, artifactID models.ArtifactID) (io.ReadCloser, error) {
	return ioutil.NopCloser(bytes.NewReader(f.data)), nil
}

func newTestDownloadContext() *JobBuildContext {
	return NewJobBuildContext(context.Background(), &documents.RunnableJob{})
}

func newTestDownloadLogger() *logging.StructuredLogger {
	return logging.NewNoOpLogPipeline().StructuredLogger().Wrap("test", "test")
}

func TestDownloadArtifact_HashMatchSucceeds(t *testing.T) {
	content := []byte("hello world, this is a test artifact")
	sum := md5.Sum(content)

	workspaceDir := t.TempDir()
	manager := NewArtifactManager(false, workspaceDir, &fakeArtifactDataAPIClient{data: content})

	artifact := &models.Artifact{
		HashType: models.HashTypeMD5,
		Hash:     hex.EncodeToString(sum[:]),
		Size:     uint64(len(content)),
	}
	artifact.Path = "some/path/artifact.bin"

	err := manager.downloadArtifact(newTestDownloadContext(), newTestDownloadLogger(), artifact)
	require.NoError(t, err)

	written, err := ioutil.ReadFile(filepath.Join(workspaceDir, artifact.Path))
	require.NoError(t, err)
	require.Equal(t, content, written)
}

func TestDownloadArtifact_HashMismatchIsRejected(t *testing.T) {
	content := []byte("hello world, this is a test artifact")

	workspaceDir := t.TempDir()
	manager := NewArtifactManager(false, workspaceDir, &fakeArtifactDataAPIClient{data: content})

	artifact := &models.Artifact{
		HashType: models.HashTypeMD5,
		Hash:     "0000000000000000000000000000000", // deliberately wrong
		Size:     uint64(len(content)),
	}
	artifact.Path = "some/path/artifact.bin"

	err := manager.downloadArtifact(newTestDownloadContext(), newTestDownloadLogger(), artifact)
	require.Error(t, err)

	// The partially-written file must not be left behind for a mismatched download.
	_, err = os.Stat(filepath.Join(workspaceDir, artifact.Path))
	require.True(t, os.IsNotExist(err), "expected corrupt/tampered artifact file to be removed")
}
