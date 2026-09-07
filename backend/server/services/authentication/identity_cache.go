package authentication

import (
	"sync"
	"time"

	"github.com/buildbeaver/buildbeaver/common/models"
)

// identityCacheTTL is how long a successful identity lookup is cached for. Disabling or
// deleting an identity can take up to this long to be reflected in already-cached lookups.
// Credential validity (e.g. Credential.IsEnabled) is checked separately, via an uncached
// credentialStore read, on every call - so revoking a credential still takes effect
// immediately regardless of this cache.
const identityCacheTTL = 10 * time.Second

type identityCacheEntry struct {
	identity  *models.Identity
	expiresAt time.Time
}

// identityCache is a small in-process cache of identityID -> Identity, used to avoid a DB round
// trip to re-read the same identity on every authenticated request (JWT, shared secret and
// client certificate auth all re-validate the identity this way). Entries are bounded by the
// number of distinct identities that actually authenticate - users, runners, service accounts -
// not by request volume, so no background eviction sweep is needed at the expected scale of a
// single BuildBeaver install; an entry is only ever cleaned up lazily, when it's next looked up
// after expiring.
type identityCache struct {
	mu      sync.Mutex
	ttl     time.Duration
	entries map[models.IdentityID]identityCacheEntry
}

func newIdentityCache(ttl time.Duration) *identityCache {
	return &identityCache{
		ttl:     ttl,
		entries: make(map[models.IdentityID]identityCacheEntry),
	}
}

func (c *identityCache) get(id models.IdentityID) (*models.Identity, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, found := c.entries[id]
	if !found {
		return nil, false
	}
	if time.Now().After(entry.expiresAt) {
		delete(c.entries, id)
		return nil, false
	}
	return entry.identity, true
}

func (c *identityCache) set(id models.IdentityID, identity *models.Identity) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[id] = identityCacheEntry{
		identity:  identity,
		expiresAt: time.Now().Add(c.ttl),
	}
}
