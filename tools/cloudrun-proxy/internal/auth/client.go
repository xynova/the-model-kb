package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/impersonate"
)

// AuthMode describes which credential path mints upstream ID tokens.
type AuthMode string

const (
	AuthModePersonalADC AuthMode = "personal_adc"
	AuthModeServiceADC  AuthMode = "service_adc"
	AuthModeImpersonate AuthMode = "impersonate"
)

// NewHTTPClient returns an http.Client that attaches Cloud Run ID tokens for audience.
// The client has no global timeout; callers should set per-request context deadlines.
func NewHTTPClient(ctx context.Context, audience, impersonateSA string) (*http.Client, AuthMode, error) {
	ts, mode, err := newIDTokenSource(ctx, audience, impersonateSA)
	if err != nil {
		return nil, "", err
	}

	client := oauth2.NewClient(ctx, ts)
	client.Timeout = 0
	return client, mode, nil
}

func newIDTokenSource(ctx context.Context, audience, impersonateSA string) (oauth2.TokenSource, AuthMode, error) {
	if audience == "" {
		return nil, "", fmt.Errorf("audience is required")
	}

	if impersonateSA != "" {
		ts, err := impersonate.IDTokenSource(ctx, impersonate.IDTokenConfig{
			Audience:        audience,
			TargetPrincipal: impersonateSA,
			IncludeEmail:    true,
		})
		if err != nil {
			return nil, "", fmt.Errorf("impersonate ID token source: %w", err)
		}
		return ts, AuthModeImpersonate, nil
	}

	if isPersonalADC() {
		ts, err := newPersonalTokenSource(audience)
		if err != nil {
			return nil, "", fmt.Errorf("personal ADC: %w", err)
		}
		return ts, AuthModePersonalADC, nil
	}

	ts, err := idtoken.NewTokenSource(ctx, audience)
	if err != nil {
		return nil, "", fmt.Errorf("service/workload ADC: %w", err)
	}
	return ts, AuthModeServiceADC, nil
}

func WithRequestTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout <= 0 {
		return parent, func() {}
	}
	return context.WithTimeout(parent, timeout)
}
