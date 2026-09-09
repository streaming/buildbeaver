package authentication

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/buildbeaver/buildbeaver/common/certificates"
	"github.com/buildbeaver/buildbeaver/common/gerror"
	"github.com/buildbeaver/buildbeaver/common/logger"
	"github.com/buildbeaver/buildbeaver/common/models"
	"github.com/buildbeaver/buildbeaver/server/store"
)

// makeSelfSignedCertData creates a minimal self-signed certificate, valid between notBefore and
// notAfter, and returns its DER-encoded bytes and public key.
func makeSelfSignedCertData(t *testing.T, notBefore, notAfter time.Time) (certificates.CertificateData, certificates.PublicKeyData) {
	t.Helper()
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{Organization: []string{"BuildBeaver Test"}},
		NotBefore:    notBefore,
		NotAfter:     notAfter,
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, pub, priv)
	require.NoError(t, err)

	publicKeyASN1, err := x509.MarshalPKIXPublicKey(pub)
	require.NoError(t, err)

	return der, publicKeyASN1
}

// fakeCredentialStore is a minimal fake store.CredentialStore that returns a fixed credential for
// a matching public key.
type fakeCredentialStore struct {
	store.CredentialStore
	publicKey certificates.PublicKeyData
	cred      *models.Credential
}

func (s *fakeCredentialStore) ReadByPublicKey(ctx context.Context, txOrNil *store.Tx, publicKey certificates.PublicKeyData) (*models.Credential, error) {
	if string(publicKey) != string(s.publicKey) {
		return nil, gerror.NewErrNotFound("Not Found")
	}
	return s.cred, nil
}

// fakeIdentityStore is a minimal fake store.IdentityStore that returns a fixed identity.
type fakeIdentityStore struct {
	store.IdentityStore
	identity *models.Identity
}

func (s *fakeIdentityStore) Read(ctx context.Context, txOrNil *store.Tx, id models.IdentityID) (*models.Identity, error) {
	if s.identity == nil || s.identity.ID != id {
		return nil, gerror.NewErrNotFound("Not Found")
	}
	return s.identity, nil
}

func newTestAuthenticationService(publicKey certificates.PublicKeyData) *AuthenticationService {
	identity := &models.Identity{ID: models.NewIdentityID()}
	cred := &models.Credential{IdentityID: identity.ID, IsEnabled: true}
	return &AuthenticationService{
		credentialStore: &fakeCredentialStore{publicKey: publicKey, cred: cred},
		identityStore:   &fakeIdentityStore{identity: identity},
		identityCache:   newIdentityCache(time.Minute),
		Log:             logger.NoOpLogFactory("test"),
	}
}

func TestAuthenticateClientCertificate_RejectsExpiredCertificate(t *testing.T) {
	now := time.Now()
	certData, publicKey := makeSelfSignedCertData(t, now.Add(-2*time.Hour), now.Add(-time.Hour))
	service := newTestAuthenticationService(publicKey)

	_, err := service.AuthenticateClientCertificate(context.Background(), certData)
	require.Error(t, err)
	require.True(t, gerror.IsUnauthorized(err))
}

func TestAuthenticateClientCertificate_RejectsNotYetValidCertificate(t *testing.T) {
	now := time.Now()
	certData, publicKey := makeSelfSignedCertData(t, now.Add(time.Hour), now.Add(2*time.Hour))
	service := newTestAuthenticationService(publicKey)

	_, err := service.AuthenticateClientCertificate(context.Background(), certData)
	require.Error(t, err)
	require.True(t, gerror.IsUnauthorized(err))
}

func TestAuthenticateClientCertificate_AcceptsValidCertificate(t *testing.T) {
	now := time.Now()
	certData, publicKey := makeSelfSignedCertData(t, now.Add(-time.Hour), now.Add(time.Hour))
	service := newTestAuthenticationService(publicKey)

	identity, err := service.AuthenticateClientCertificate(context.Background(), certData)
	require.NoError(t, err)
	require.NotNil(t, identity)
}
