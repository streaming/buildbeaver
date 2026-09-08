package github

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/buildbeaver/buildbeaver/common/gerror"
	"github.com/buildbeaver/buildbeaver/common/logger"
)

const testWebhookPayload = `{"action":"added_to_repository"}`

func sha256Signature(t *testing.T, secret []byte, payload []byte) string {
	t.Helper()
	mac := hmac.New(sha256.New, secret)
	_, err := mac.Write(payload)
	require.NoError(t, err)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

// newTestGitHubService returns a GitHubService with the given webhook secret configured, and no
// other dependencies wired up. Only usable for tests that exercise webhook signature verification
// against event types (such as "team_add") whose handling doesn't touch any other service or store.
func newTestGitHubService(webhookSecret []byte) *GitHubService {
	return &GitHubService{
		config: AppConfig{WebhookSecret: webhookSecret},
		Log:    logger.NoOpLogFactory("test"),
	}
}

func TestHandleWebhookEvent_ValidSignatureAccepted(t *testing.T) {
	secret := []byte("test-secret")
	payload := []byte(testWebhookPayload)
	s := newTestGitHubService(secret)

	event := &WebhookEvent{
		EventType:    "team_add",
		Signature256: sha256Signature(t, secret, payload),
		Payload:      bytes.NewReader(payload),
	}
	err := s.HandleWebhookEvent(context.Background(), event)
	require.NoError(t, err)
}

func TestHandleWebhookEvent_InvalidSignatureRejected(t *testing.T) {
	secret := []byte("test-secret")
	payload := []byte(testWebhookPayload)
	s := newTestGitHubService(secret)

	event := &WebhookEvent{
		EventType:    "team_add",
		Signature256: sha256Signature(t, []byte("wrong-secret"), payload),
		Payload:      bytes.NewReader(payload),
	}
	err := s.HandleWebhookEvent(context.Background(), event)
	require.Error(t, err)
	require.True(t, gerror.IsUnauthorized(err))
}

func TestHandleWebhookEvent_TamperedPayloadRejected(t *testing.T) {
	secret := []byte("test-secret")
	payload := []byte(testWebhookPayload)
	s := newTestGitHubService(secret)

	event := &WebhookEvent{
		EventType:    "team_add",
		Signature256: sha256Signature(t, secret, payload),
		Payload:      bytes.NewReader([]byte(`{"action":"removed_from_repository"}`)),
	}
	err := s.HandleWebhookEvent(context.Background(), event)
	require.Error(t, err)
	require.True(t, gerror.IsUnauthorized(err))
}

func TestHandleWebhookEvent_NoSecretConfiguredFailsClosed(t *testing.T) {
	payload := []byte(testWebhookPayload)
	s := newTestGitHubService(nil)

	// Even a correctly-formed signature (against some secret the attacker knows) must be
	// rejected when the server has no secret configured to check it against.
	event := &WebhookEvent{
		EventType:    "team_add",
		Signature256: sha256Signature(t, []byte("whatever"), payload),
		Payload:      bytes.NewReader(payload),
	}
	err := s.HandleWebhookEvent(context.Background(), event)
	require.Error(t, err)
	require.True(t, gerror.IsUnauthorized(err))
}

func TestWebhookHandler_MissingSignatureHeaderRejected(t *testing.T) {
	s := newTestGitHubService([]byte("test-secret"))
	handler, err := s.WebhookHandler()
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(testWebhookPayload)))
	req.Header.Set("X-GitHub-Event", "team_add")
	// Deliberately no X-Hub-Signature-256 header set.
	rec := httptest.NewRecorder()

	handler(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestWebhookHandler_InvalidSignatureReturns401(t *testing.T) {
	s := newTestGitHubService([]byte("test-secret"))
	handler, err := s.WebhookHandler()
	require.NoError(t, err)

	payload := []byte(testWebhookPayload)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(payload))
	req.Header.Set("X-GitHub-Event", "team_add")
	req.Header.Set("X-Hub-Signature-256", sha256Signature(t, []byte("wrong-secret"), payload))
	rec := httptest.NewRecorder()

	handler(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestWebhookHandler_ValidSignatureReturns200(t *testing.T) {
	secret := []byte("test-secret")
	s := newTestGitHubService(secret)
	handler, err := s.WebhookHandler()
	require.NoError(t, err)

	payload := []byte(testWebhookPayload)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(payload))
	req.Header.Set("X-GitHub-Event", "team_add")
	req.Header.Set("X-Hub-Signature-256", sha256Signature(t, secret, payload))
	rec := httptest.NewRecorder()

	handler(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}
