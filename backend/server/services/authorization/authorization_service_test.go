package authorization

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/buildbeaver/buildbeaver/common/logger"
	"github.com/buildbeaver/buildbeaver/common/models"
	"github.com/buildbeaver/buildbeaver/server/store"
)

// countingAuthorizationStore is a minimal fake store.AuthorizationStore that counts calls to
// CountGrantsForOperation, so tests can assert on whether the request cache avoided a query.
type countingAuthorizationStore struct {
	count     int
	callCount int
}

func (s *countingAuthorizationStore) CountGrantsForOperation(
	ctx context.Context,
	txOrNil *store.Tx,
	identityID models.IdentityID,
	operation *models.Operation,
	resourceID models.ResourceID,
) (int, error) {
	s.callCount++
	return s.count, nil
}

func TestIsAuthorizedUsesRequestCacheWhenPresent(t *testing.T) {
	fakeStore := &countingAuthorizationStore{count: 1} // 1 grant found -> authorized
	service := &AuthorizationService{authorizationStore: fakeStore, Log: logger.NewNoOpLog()}

	identityID := models.NewIdentityID()
	operation := &models.Operation{Name: "read", ResourceKind: "repo"}
	resourceID := models.NewResourceID("repo")

	ctx := WithRequestCache(context.Background())

	allowed, err := service.IsAuthorized(ctx, identityID, operation, resourceID)
	require.NoError(t, err)
	require.True(t, allowed)
	require.Equal(t, 1, fakeStore.callCount, "first check should query the store")

	allowed, err = service.IsAuthorized(ctx, identityID, operation, resourceID)
	require.NoError(t, err)
	require.True(t, allowed)
	require.Equal(t, 1, fakeStore.callCount, "second check with the same request cache should be served from cache")
}

func TestIsAuthorizedQueriesEveryTimeWithoutRequestCache(t *testing.T) {
	fakeStore := &countingAuthorizationStore{count: 1}
	service := &AuthorizationService{authorizationStore: fakeStore, Log: logger.NewNoOpLog()}

	identityID := models.NewIdentityID()
	operation := &models.Operation{Name: "read", ResourceKind: "repo"}
	resourceID := models.NewResourceID("repo")

	ctx := context.Background() // no request cache attached

	_, err := service.IsAuthorized(ctx, identityID, operation, resourceID)
	require.NoError(t, err)
	_, err = service.IsAuthorized(ctx, identityID, operation, resourceID)
	require.NoError(t, err)

	require.Equal(t, 2, fakeStore.callCount, "without a request cache, every call should query the store")
}

func TestIsAuthorizedRequestCacheDistinguishesResources(t *testing.T) {
	fakeStore := &countingAuthorizationStore{count: 0} // no grant found -> not authorized
	service := &AuthorizationService{authorizationStore: fakeStore, Log: logger.NewNoOpLog()}

	identityID := models.NewIdentityID()
	operation := &models.Operation{Name: "read", ResourceKind: "repo"}

	ctx := WithRequestCache(context.Background())

	allowed, err := service.IsAuthorized(ctx, identityID, operation, models.NewResourceID("repo"))
	require.NoError(t, err)
	require.False(t, allowed)
	require.Equal(t, 1, fakeStore.callCount)

	// A different resource ID must not be served from the first resource's cache entry.
	allowed, err = service.IsAuthorized(ctx, identityID, operation, models.NewResourceID("repo"))
	require.NoError(t, err)
	require.False(t, allowed)
	require.Equal(t, 2, fakeStore.callCount, "a different resource id should not hit the cache")
}
