package secret_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/buildbeaver/buildbeaver/common/models"
	"github.com/buildbeaver/buildbeaver/server/app/server_test"
)

// TestListPlaintextByRepoIDAndNames confirms that ListPlaintextByRepoIDAndNames returns only the
// secrets whose plaintext key was requested, not every secret in the repo - this is the core of
// the fix for the runner previously loading every one of a repo's secrets into memory regardless
// of whether the running build referenced them.
func TestListPlaintextByRepoIDAndNames(t *testing.T) {
	ctx := context.Background()
	app, cleanup, err := server_test.New(server_test.TestConfig(t))
	require.NoError(t, err)
	defer cleanup()

	company := server_test.CreateCompanyLegalEntity(t, ctx, app, "", "", "")
	repo := server_test.CreateRepo(t, ctx, app, company.ID)

	_, err = app.SecretService.Create(ctx, nil, repo.ID, "secret-a", "value-a", false)
	require.NoError(t, err)
	_, err = app.SecretService.Create(ctx, nil, repo.ID, "secret-b", "value-b", false)
	require.NoError(t, err)
	_, err = app.SecretService.Create(ctx, nil, repo.ID, "secret-c", "value-c", false)
	require.NoError(t, err)

	pagination := models.NewPagination(10, nil)

	t.Run("returns only the requested subset", func(t *testing.T) {
		secrets, _, err := app.SecretService.ListPlaintextByRepoIDAndNames(ctx, nil, repo.ID, []string{"secret-a", "secret-c"}, pagination)
		require.NoError(t, err)
		require.Len(t, secrets, 2)
		byKey := map[string]string{}
		for _, s := range secrets {
			byKey[s.Key] = s.Value
		}
		require.Equal(t, "value-a", byKey["secret-a"])
		require.Equal(t, "value-c", byKey["secret-c"])
		require.NotContains(t, byKey, "secret-b")
	})

	t.Run("empty names returns no secrets, not every secret", func(t *testing.T) {
		secrets, _, err := app.SecretService.ListPlaintextByRepoIDAndNames(ctx, nil, repo.ID, nil, pagination)
		require.NoError(t, err)
		require.Empty(t, secrets)
	})

	t.Run("unrecognized name returns no matching secret", func(t *testing.T) {
		secrets, _, err := app.SecretService.ListPlaintextByRepoIDAndNames(ctx, nil, repo.ID, []string{"does-not-exist"}, pagination)
		require.NoError(t, err)
		require.Empty(t, secrets)
	})
}
