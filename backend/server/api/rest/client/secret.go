package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/buildbeaver/buildbeaver/common/models"
	"github.com/buildbeaver/buildbeaver/server/api/rest/documents"
)

type paginatedSecretResponse struct {
	*documents.PaginatedResponse
	Results []*models.SecretPlaintext `json:"results"` // TODO these should be documents
}

// GetSecretsPlaintextByNames gets, in plaintext, only the secrets for the specified repo whose
// plaintext key is one of names. Only the subset of a repo's secrets actually needed should be
// requested - see ArtifactManager and SecretStore for why this matters.
func (a *APIClient) GetSecretsPlaintextByNames(ctx context.Context, repoID models.RepoID, names []string) ([]*models.SecretPlaintext, error) {
	url := fmt.Sprintf("/api/v1/runner/repos/%s/secrets/search", repoID)
	req := &documents.SecretSearchRequest{Names: names}
	code, _, body, err := a.post(ctx, nil, url, req)
	if err != nil {
		return nil, err
	}
	if !a.isOneOf(code, []int{http.StatusOK}) {
		return nil, a.makeHTTPError(code, body)
	}
	doc := &paginatedSecretResponse{}
	err = json.Unmarshal(body, doc)
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing response body: %s", string(body[:]))
	}
	return doc.Results, nil
}
