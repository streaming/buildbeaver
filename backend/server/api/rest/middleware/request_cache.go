package middleware

import (
	"net/http"

	"github.com/buildbeaver/buildbeaver/server/services/authorization"
)

// MakeAuthorizationRequestCache makes a middleware that attaches a fresh, empty per-request
// authorization cache to the request context, so that AuthorizationService.IsAuthorized can
// memoize results for the lifetime of a single request. Should be applied to every route that
// might call IsAuthorized.
func MakeAuthorizationRequestCache() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := authorization.WithRequestCache(r.Context())
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
