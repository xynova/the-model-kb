package auth

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	"golang.org/x/oauth2"
)

// personalTokenSource mints audience-bound ID tokens for local dev.
//
// google.golang.org/api/idtoken does not support authorized_user ADC, so we use
// gcloud — the same audience-scoped flow as:
//
//	gcloud auth print-identity-token --audiences=$AUDIENCE
//
// Requires gcloud auth application-default login (ADC file) and an active
// gcloud user session (gcloud auth login) for the same account.
type personalTokenSource struct {
	audience string
	mu       sync.Mutex
	tok      *oauth2.Token
}

func newPersonalTokenSource(audience string) (oauth2.TokenSource, error) {
	if audience == "" {
		return nil, fmt.Errorf("audience is required")
	}
	if _, err := exec.LookPath("gcloud"); err != nil {
		return nil, fmt.Errorf("gcloud not found in PATH; install Google Cloud SDK or set IMPERSONATE_SERVICE_ACCOUNT")
	}
	if !isPersonalADC() {
		return nil, fmt.Errorf("personal ADC not found; run: gcloud auth application-default login")
	}
	return &personalTokenSource{audience: audience}, nil
}

func (s *personalTokenSource) Token() (*oauth2.Token, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.tok != nil && s.tok.Valid() {
		return s.tok, nil
	}

	cmd := exec.Command(
		"gcloud", "auth", "print-identity-token",
		"--audiences="+s.audience,
	)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	out, err := cmd.Output()
	if err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return nil, fmt.Errorf("gcloud auth print-identity-token: %s (re-auth: gcloud auth login && gcloud auth application-default login)", msg)
	}

	token := strings.TrimSpace(string(out))
	if token == "" {
		return nil, fmt.Errorf("gcloud auth print-identity-token returned empty token")
	}

	expiry := time.Now().Add(55 * time.Minute)
	if exp, ok := jwtExpiry(token); ok {
		expiry = exp.Add(-time.Minute)
	}

	s.tok = &oauth2.Token{
		AccessToken: token,
		TokenType:   "Bearer",
		Expiry:      expiry,
	}
	return s.tok, nil
}
