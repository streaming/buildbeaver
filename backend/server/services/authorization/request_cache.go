package authorization

import (
	"context"
	"sync"

	"github.com/buildbeaver/buildbeaver/common/models"
)

// requestCacheKey identifies a single IsAuthorized question: can this identity perform this
// operation on this resource.
type requestCacheKey struct {
	identityID models.IdentityID
	operation  models.Operation
	resourceID models.ResourceID
}

// RequestCache memoizes IsAuthorized results for the lifetime of a single request, so that a
// handler which checks the same identity/operation/resource combination more than once (for
// example AuthorizedResourceID followed by a related permission check) only pays for one
// authorization query. It carries no state beyond a single request and is never persisted, so
// it introduces no staleness risk.
type RequestCache struct {
	mu      sync.Mutex
	results map[requestCacheKey]bool
}

type requestCacheContextKey struct{}

// WithRequestCache returns a context with a fresh, empty per-request authorization cache
// attached. Called once per incoming HTTP request (see server/api/rest/middleware) before any
// authorization checks are made against it; IsAuthorized falls back to its uncached behaviour
// if no cache has been attached to the context.
func WithRequestCache(ctx context.Context) context.Context {
	return context.WithValue(ctx, requestCacheContextKey{}, &RequestCache{
		results: make(map[requestCacheKey]bool),
	})
}

func requestCacheFromContext(ctx context.Context) *RequestCache {
	cache, _ := ctx.Value(requestCacheContextKey{}).(*RequestCache)
	return cache
}

func (c *RequestCache) get(key requestCacheKey) (allowed bool, found bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	allowed, found = c.results[key]
	return allowed, found
}

func (c *RequestCache) set(key requestCacheKey, allowed bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.results[key] = allowed
}
