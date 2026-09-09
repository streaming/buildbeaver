package certificates

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// makeSelfSignedCertData creates a minimal self-signed certificate, valid between notBefore and
// notAfter, and returns its DER-encoded bytes.
func makeSelfSignedCertData(t *testing.T, notBefore, notAfter time.Time) CertificateData {
	t.Helper()
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{Organization: []string{"BuildBeaver Test"}},
		NotBefore:    notBefore,
		NotAfter:     notAfter,
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, publicKey, privateKey)
	require.NoError(t, err)
	return der
}

func TestCheckCertificateValidityPeriod(t *testing.T) {
	now := time.Now()

	t.Run("valid certificate", func(t *testing.T) {
		certData := makeSelfSignedCertData(t, now.Add(-time.Hour), now.Add(time.Hour))
		err := CheckCertificateValidityPeriod(certData, now)
		require.NoError(t, err)
	})

	t.Run("expired certificate", func(t *testing.T) {
		certData := makeSelfSignedCertData(t, now.Add(-2*time.Hour), now.Add(-time.Hour))
		err := CheckCertificateValidityPeriod(certData, now)
		require.Error(t, err)
	})

	t.Run("not yet valid certificate", func(t *testing.T) {
		certData := makeSelfSignedCertData(t, now.Add(time.Hour), now.Add(2*time.Hour))
		err := CheckCertificateValidityPeriod(certData, now)
		require.Error(t, err)
	})

	t.Run("malformed certificate data", func(t *testing.T) {
		err := CheckCertificateValidityPeriod([]byte("not a certificate"), now)
		require.Error(t, err)
	})
}
