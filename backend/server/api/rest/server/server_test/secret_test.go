package api_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/buildbeaver/buildbeaver/server/api/rest/client/clienttest"
	"github.com/buildbeaver/buildbeaver/server/app/server_test"
)

// TestGetSecretsPlaintextByNames_OnlyReturnsRequestedSecrets exercises the real runner-facing
// endpoint end-to-end (client -> HTTP -> handler -> service -> store), confirming that a runner
// requesting a specific subset of a repo's secrets by name gets back only that subset - not every
// secret in the repo, which is what closed the excess-blast-radius gap this endpoint replaces.
func TestGetSecretsPlaintextByNames_OnlyReturnsRequestedSecrets(t *testing.T) {
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

	_, err = app.SecretService.Create(ctx, nil, repo.ID, "wanted-secret", "wanted-value", false)
	require.NoError(t, err)
	_, err = app.SecretService.Create(ctx, nil, repo.ID, "unwanted-secret", "unwanted-value", false)
	require.NoError(t, err)

	secrets, err := apiClient.GetSecretsPlaintextByNames(ctx, repo.ID, []string{"wanted-secret"})
	require.NoError(t, err)
	require.Len(t, secrets, 1)
	require.Equal(t, "wanted-secret", secrets[0].Key)
	require.Equal(t, "wanted-value", secrets[0].Value)
}
