package authentication

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/buildbeaver/buildbeaver/common/gerror"
	"github.com/buildbeaver/buildbeaver/common/models"
	"github.com/buildbeaver/buildbeaver/server/store"
)

func TestIdentityCacheGetSet(t *testing.T) {
	cache := newIdentityCache(time.Hour)
	id := models.NewIdentityID()
	identity := &models.Identity{ID: id}

	_, found := cache.get(id)
	require.False(t, found, "unset key should not be found")

	cache.set(id, identity)
	got, found := cache.get(id)
	require.True(t, found)
	require.Same(t, identity, got)

	otherID := models.NewIdentityID()
	_, found = cache.get(otherID)
	require.False(t, found, "a different key should not be found")
}

func TestIdentityCacheExpiry(t *testing.T) {
	cache := newIdentityCache(10 * time.Millisecond)
	id := models.NewIdentityID()
	cache.set(id, &models.Identity{ID: id})

	_, found := cache.get(id)
	require.True(t, found, "entry should be found immediately after being set")

	time.Sleep(20 * time.Millisecond)

	_, found = cache.get(id)
	require.False(t, found, "entry should no longer be found once its TTL has passed")
}

// countingIdentityStore is a minimal fake store.IdentityStore that counts calls to Read, so
// tests can assert on whether the identity cache avoided a store read.
type countingIdentityStore struct {
	store.IdentityStore
	identity  *models.Identity
	readCalls int
}

func (s *countingIdentityStore) Read(ctx context.Context, txOrNil *store.Tx, id models.IdentityID) (*models.Identity, error) {
	s.readCalls++
	if s.identity == nil || s.identity.ID != id {
		return nil, gerror.NewErrNotFound("Not Found")
	}
	return s.identity, nil
}

func TestAuthenticationServiceReadIdentityUsesCache(t *testing.T) {
	identity := &models.Identity{ID: models.NewIdentityID()}
	fakeStore := &countingIdentityStore{identity: identity}
	service := &AuthenticationService{
		identityStore: fakeStore,
		identityCache: newIdentityCache(time.Hour),
	}
	ctx := context.Background()

	got, err := service.readIdentity(ctx, identity.ID)
	require.NoError(t, err)
	require.Same(t, identity, got)
	require.Equal(t, 1, fakeStore.readCalls, "first call should read through to the store")

	got, err = service.readIdentity(ctx, identity.ID)
	require.NoError(t, err)
	require.Same(t, identity, got)
	require.Equal(t, 1, fakeStore.readCalls, "second call within the TTL should be served from cache")
}
