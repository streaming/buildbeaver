package authorization

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/buildbeaver/buildbeaver/common/models"
)

func TestRequestCacheFromContextWithNoCacheAttached(t *testing.T) {
	require.Nil(t, requestCacheFromContext(context.Background()),
		"a context with no cache attached should yield a nil cache, not panic or a usable cache")
}

func TestRequestCacheGetSet(t *testing.T) {
	ctx := WithRequestCache(context.Background())
	cache := requestCacheFromContext(ctx)
	require.NotNil(t, cache)

	key := requestCacheKey{
		identityID: models.NewIdentityID(),
		operation:  models.Operation{Name: "read", ResourceKind: "repo"},
		resourceID: models.NewResourceID("repo"),
	}

	_, found := cache.get(key)
	require.False(t, found, "unset key should not be found")

	cache.set(key, true)
	allowed, found := cache.get(key)
	require.True(t, found)
	require.True(t, allowed)

	cache.set(key, false)
	allowed, found = cache.get(key)
	require.True(t, found)
	require.False(t, allowed, "setting the same key again should overwrite the previous value")
}

func TestRequestCacheIsolatedByKey(t *testing.T) {
	ctx := WithRequestCache(context.Background())
	cache := requestCacheFromContext(ctx)

	identityA := models.NewIdentityID()
	identityB := models.NewIdentityID()
	resource := models.NewResourceID("repo")
	operation := models.Operation{Name: "read", ResourceKind: "repo"}

	cache.set(requestCacheKey{identityID: identityA, operation: operation, resourceID: resource}, true)

	_, found := cache.get(requestCacheKey{identityID: identityB, operation: operation, resourceID: resource})
	require.False(t, found, "a different identity should not share a cache entry")

	_, found = cache.get(requestCacheKey{identityID: identityA, operation: models.Operation{Name: "update", ResourceKind: "repo"}, resourceID: resource})
	require.False(t, found, "a different operation should not share a cache entry")
}

func TestRequestCacheIsolatedBetweenContexts(t *testing.T) {
	key := requestCacheKey{
		identityID: models.NewIdentityID(),
		operation:  models.Operation{Name: "read", ResourceKind: "repo"},
		resourceID: models.NewResourceID("repo"),
	}

	firstCtx := WithRequestCache(context.Background())
	requestCacheFromContext(firstCtx).set(key, true)

	secondCtx := WithRequestCache(context.Background())
	_, found := requestCacheFromContext(secondCtx).get(key)
	require.False(t, found, "a fresh request cache should not see entries set on another one")
}
