package client

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"

	"github.com/buildbeaver/buildbeaver/common/models"
	"github.com/buildbeaver/buildbeaver/server/api/rest/documents"
)

type artifactSearchPaginator struct {
	*searchPaginator
}

func newArtifactSearchPaginator(
	apiClient *APIClient,
	url string,
	search *documents.ArtifactSearchRequest) *artifactSearchPaginator {
	return &artifactSearchPaginator{
		searchPaginator: newSearchPaginator(apiClient, url, search),
	}
}

func (a *artifactSearchPaginator) Next(ctx context.Context) ([]*models.Artifact, error) {
	raw, err := a.next(ctx)
	if err != nil {
		return nil, err
	}
	var results []*models.Artifact
	err = json.Unmarshal(raw, &results)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling artifacts: %w", err)
	}
	return results, nil
}

// CreateArtifact registers a new artifact against the specified step.
func (a *APIClient) CreateArtifact(
	ctx context.Context,
	jobID models.JobID,
	groupName models.ResourceName,
	relativePath string,
	reader io.ReadSeeker) (*documents.Artifact, error) {

	// Hash the artifact data before uploading so the server can verify it arrived intact, then
	// seek back to the start so postStream (via retryablehttp) can read the same data as the
	// request body. The server compares this as a hex digest (see ArtifactService.Create), not
	// the base64 encoding RFC 1864 specifies for a Content-MD5 header.
	hasher := md5.New()
	_, err := io.Copy(hasher, reader)
	if err != nil {
		return nil, fmt.Errorf("error hashing artifact data: %w", err)
	}
	_, err = reader.Seek(0, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("error seeking back to start of artifact data: %w", err)
	}

	url := fmt.Sprintf("/api/v1/runner/jobs/%s/artifacts", jobID)
	headers := http.Header{
		"X-BuildBeaver-Artifact-Path":  []string{relativePath},
		"X-BuildBeaver-Artifact-Group": []string{groupName.String()},
		"Content-MD5":                  []string{hex.EncodeToString(hasher.Sum(nil))},
	}
	code, headers, body, err := a.postStream(ctx, headers, url, reader)
	if err != nil {
		return nil, fmt.Errorf("error in request: %w", err)
	}
	defer body.Close()
	buf, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("error reading body: %w", err)
	}
	if !a.isOneOf(code, []int{http.StatusOK, http.StatusCreated, http.StatusNoContent}) {
		return nil, a.makeHTTPError(code, buf)
	}
	resDoc := &documents.Artifact{}
	err = json.Unmarshal(buf, resDoc)
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing response body: %s", string(buf[:]))
	}
	return resDoc, nil
}

// GetArtifactData returns a reader to the data of an artifact.
// It is the callers responsibility to close the reader.
func (a *APIClient) GetArtifactData(ctx context.Context, artifactID models.ArtifactID) (io.ReadCloser, error) {
	url := fmt.Sprintf("/api/v1/runner/artifacts/%s/data", artifactID)
	code, _, body, err := a.getStream(ctx, nil, url)
	if err != nil {
		return nil, err
	}
	if !a.isOneOf(code, []int{http.StatusOK}) {
		body.Close()
		return nil, a.makeHTTPError(code, nil)
	}
	return body, nil
}

// SearchArtifacts searches all artifacts for a build. Use pager to page through results, if any.
func (a *APIClient) SearchArtifacts(ctx context.Context, buildID models.BuildID, search *models.ArtifactSearch) (models.ArtifactSearchPaginator, error) {
	doc := &documents.ArtifactSearchRequest{
		ArtifactSearch: search,
	}
	url := fmt.Sprintf("/api/v1/runner/builds/%s/artifacts/search", buildID)
	paginator := newArtifactSearchPaginator(a, url, doc)
	return paginator, nil
}
