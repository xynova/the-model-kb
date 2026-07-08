package upstream

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/xynova/cloudrun-proxy/internal/auth"
	"github.com/xynova/cloudrun-proxy/internal/config"
	"github.com/xynova/cloudrun-proxy/internal/openai"
)

const maxBodyBytes = 32 << 20 // 32 MiB

var retryableStatuses = map[int]bool{
	http.StatusBadGateway:         true,
	http.StatusServiceUnavailable: true,
	http.StatusGatewayTimeout:     true,
}

type Router struct {
	cfg             config.Config
	models          map[string]modelRoute
	slots           chan struct{}
	logger          *slog.Logger
	rateLimitMu     sync.Mutex
	rateLimitUntil  time.Time
}

type modelRoute struct {
	config.ResolvedModel
	httpClient *http.Client
	authMode   auth.AuthMode
}

func NewRouter(ctx context.Context, cfg config.Config, logger *slog.Logger) (*Router, error) {
	clientCache := &authClientCache{logger: logger}
	models := make(map[string]modelRoute, len(cfg.Models))

	for id := range cfg.Models {
		resolved, ok := cfg.Resolve(id)
		if !ok {
			continue
		}
		client, mode, err := clientCache.get(ctx, resolved.Audience, resolved.ImpersonateSA)
		if err != nil {
			return nil, fmt.Errorf("models.%s auth: %w", id, err)
		}
		models[id] = modelRoute{
			ResolvedModel: resolved,
			httpClient:    client,
			authMode:      mode,
		}
		logger.Info("model route",
			"model", id,
			"audience", resolved.Audience,
			"url", resolved.URL,
			"upstream_model", resolved.UpstreamModel,
			"impersonate", resolved.ImpersonateSA,
			"auth_mode", mode,
			"timeout", resolved.RequestTimeout.String(),
		)
	}

	n := cfg.Server.MaxConcurrent
	if n <= 0 {
		n = 1
	}
	logger.Info("upstream concurrency",
		"max_concurrent", n,
		"reject_when_busy", cfg.Server.RejectWhenBusy,
		"rate_limit_cooldown", cfg.Server.RateLimitCooldown.String(),
	)

	return &Router{
		cfg:    cfg,
		models: models,
		slots:  make(chan struct{}, n),
		logger: logger,
	}, nil
}

type authClientCache struct {
	mu      sync.Mutex
	clients map[string]cachedClient
	logger  *slog.Logger
}

type cachedClient struct {
	client *http.Client
	mode   auth.AuthMode
}

func (c *authClientCache) get(ctx context.Context, audience, impersonateSA string) (*http.Client, auth.AuthMode, error) {
	key := audience + "|" + impersonateSA

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.clients == nil {
		c.clients = make(map[string]cachedClient)
	}
	if entry, ok := c.clients[key]; ok {
		return entry.client, entry.mode, nil
	}

	client, mode, err := auth.NewHTTPClient(ctx, audience, impersonateSA)
	if err != nil {
		return nil, "", err
	}
	c.clients[key] = cachedClient{client: client, mode: mode}
	if c.logger != nil {
		c.logger.Debug("auth client created",
			"audience", audience,
			"impersonate", impersonateSA,
			"auth_mode", mode,
		)
	}
	return client, mode, nil
}

func (r *Router) ChatCompletions(w http.ResponseWriter, req *http.Request) {
	start := time.Now()

	if err := openai.CheckLocalSecret(req, r.cfg.Server.LocalSecret); err != nil {
		r.logger.Warn("local auth rejected", "error", err)
		openai.WriteError(w, http.StatusUnauthorized, "authentication_error", err.Error())
		return
	}

	body, err := openai.ReadBody(req, maxBodyBytes)
	if err != nil {
		openai.WriteError(w, http.StatusBadRequest, "invalid_request_error", err.Error())
		return
	}

	modelID, err := openai.RequestModel(body)
	if err != nil {
		openai.WriteError(w, http.StatusBadRequest, "invalid_request_error", err.Error())
		return
	}

	route, ok := r.models[modelID]
	if !ok {
		r.logger.Warn("unknown model",
			"model", modelID,
			"available", strings.Join(r.cfg.ModelIDs(), ", "),
		)
		openai.WriteError(w, http.StatusBadRequest, "invalid_request_error",
			fmt.Sprintf("unknown model %q (available: %s)", modelID, strings.Join(r.cfg.ModelIDs(), ", ")))
		return
	}

	if err := openai.ValidateChatRequest(body); err != nil {
		openai.WriteError(w, http.StatusBadRequest, "invalid_request_error", err.Error())
		return
	}

	clientSummary := openai.SummarizeChatRequest(body)
	r.logger.Debug("client request", "summary", clientSummary)

	body, err = openai.PrepareUpstreamChatBody(body, route.ResolvedModel)
	if err != nil {
		openai.WriteError(w, http.StatusBadRequest, "invalid_request_error", err.Error())
		return
	}

	upstreamSummary := openai.SummarizeUpstreamRequest(body, route.UpstreamModel)
	r.logger.Debug("upstream request prepared", "summary", upstreamSummary)

	streaming := openai.IsStreaming(body)
	upstreamURL := route.URL + route.ChatPath

	if retryAfter := r.rateLimitRetryAfter(); retryAfter > 0 {
		r.logger.Warn("rate limit cooldown active",
			"model", modelID,
			"retry_after_s", retryAfter,
		)
		openai.WriteErrorWithRetryAfter(w, http.StatusTooManyRequests, "rate_limit_error",
			"Cloud Run rate limit cooldown — wait before retrying", retryAfter)
		return
	}

	if !r.tryAcquireUpstream() {
		r.logger.Warn("rejected busy",
			"model", modelID,
			"max_concurrent", r.cfg.Server.MaxConcurrent,
		)
		openai.WriteErrorWithRetryAfter(w, http.StatusTooManyRequests, "rate_limit_error",
			fmt.Sprintf("proxy upstream busy (max_concurrent=%d) — wait for the current request to finish", r.cfg.Server.MaxConcurrent),
			5)
		return
	}
	defer r.releaseUpstream()

	r.logger.Info("forwarding",
		"model", modelID,
		"upstream", upstreamURL,
		"upstream_model", route.UpstreamModel,
		"auth_mode", route.authMode,
		"stream", streaming,
		"bytes", len(body),
	)

	upstreamStart := time.Now()
	var resp *http.Response
	var upstreamCancel context.CancelFunc
	if streaming {
		resp, upstreamCancel, err = r.doOnce(req, route, upstreamURL, body)
	} else {
		resp, upstreamCancel, err = r.doWithRetries(req, route, upstreamURL, body, modelID)
	}
	upstreamDuration := time.Since(upstreamStart)

	if err != nil {
		logArgs := []any{
			"model", modelID,
			"upstream", upstreamURL,
			"auth_mode", route.authMode,
			"upstream_ms", upstreamDuration.Milliseconds(),
			"total_ms", time.Since(start).Milliseconds(),
			"error", err,
		}
		if hint := auth.ErrorHint(err); hint != "" {
			logArgs = append(logArgs, "hint", hint)
		}
		r.logger.Error("upstream failed", logArgs...)
		openai.WriteError(w, http.StatusBadGateway, "upstream_error", err.Error())
		return
	}
	defer resp.Body.Close()
	if upstreamCancel != nil {
		defer upstreamCancel()
	}

	r.logger.Info("upstream responded",
		"model", modelID,
		"status", resp.StatusCode,
		"stream", streaming,
		"upstream_ms", upstreamDuration.Milliseconds(),
	)

	if resp.StatusCode == http.StatusTooManyRequests {
		r.markRateLimited()
		r.logger.Warn("upstream rate limited",
			"model", modelID,
			"cooldown", r.cfg.Server.RateLimitCooldown.String(),
			"hint", "client retried too fast — cooldown applied",
		)
	}

	copyResponse(w, resp, streaming, r.logger, modelID, req.Context())
	r.logger.Info("completed",
		"model", modelID,
		"total_ms", time.Since(start).Milliseconds(),
	)
}

func (r *Router) tryAcquireUpstream() bool {
	if r.cfg.Server.RejectWhenBusy {
		select {
		case r.slots <- struct{}{}:
			return true
		default:
			return false
		}
	}
	r.slots <- struct{}{}
	return true
}

func (r *Router) releaseUpstream() {
	<-r.slots
}

func (r *Router) markRateLimited() {
	cooldown := r.cfg.Server.RateLimitCooldown
	if cooldown <= 0 {
		cooldown = 60 * time.Second
	}
	until := time.Now().Add(cooldown)
	r.rateLimitMu.Lock()
	if until.After(r.rateLimitUntil) {
		r.rateLimitUntil = until
	}
	r.rateLimitMu.Unlock()
}

func (r *Router) rateLimitRetryAfter() int {
	r.rateLimitMu.Lock()
	until := r.rateLimitUntil
	r.rateLimitMu.Unlock()
	if time.Now().After(until) {
		return 0
	}
	sec := int(time.Until(until).Seconds())
	if sec < 1 {
		return 1
	}
	return sec
}

func (r *Router) doWithRetries(req *http.Request, route modelRoute, url string, body []byte, modelID string) (*http.Response, context.CancelFunc, error) {
	const maxAttempts = 3
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		resp, cancel, err := r.doOnce(req, route, url, body)
		if err != nil {
			lastErr = err
			if attempt < maxAttempts {
				r.logger.Warn("upstream attempt failed",
					"model", modelID,
					"attempt", attempt,
					"error", err,
				)
				time.Sleep(retryDelay(attempt))
				continue
			}
			return nil, nil, err
		}

		if !retryableStatuses[resp.StatusCode] || attempt == maxAttempts {
			return resp, cancel, nil
		}

		r.logger.Warn("retryable upstream status",
			"model", modelID,
			"attempt", attempt,
			"status", resp.StatusCode,
		)
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		cancel()
		time.Sleep(retryDelay(attempt))
	}

	return nil, nil, lastErr
}

func (r *Router) doOnce(req *http.Request, route modelRoute, url string, body []byte) (*http.Response, context.CancelFunc, error) {
	// Do not cancel upstream when the local client disconnects (RooCode often
	// times out before Cloud Run returns the first token). Only apply the
	// configured upstream timeout. The cancel func must outlive doOnce — callers
	// defer it until the response body is fully read.
	ctx, cancel := auth.WithRequestTimeout(context.WithoutCancel(req.Context()), route.RequestTimeout)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		cancel()
		return nil, nil, fmt.Errorf("build upstream request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if accept := req.Header.Get("Accept"); accept != "" {
		httpReq.Header.Set("Accept", accept)
	}

	r.logger.Debug("upstream http post",
		"url", url,
		"timeout", route.RequestTimeout.String(),
		"bytes", len(body),
	)

	resp, err := route.httpClient.Do(httpReq)
	if err != nil {
		cancel()
		return nil, nil, err
	}
	return resp, cancel, nil
}

func copyResponse(w http.ResponseWriter, resp *http.Response, streaming bool, logger *slog.Logger, modelID string, clientCtx context.Context) {
	for k, vals := range resp.Header {
		lk := strings.ToLower(k)
		if lk == "transfer-encoding" || lk == "connection" || lk == "content-length" {
			continue
		}
		for _, v := range vals {
			w.Header().Add(k, v)
		}
	}

	if resp.StatusCode >= 400 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			openai.WriteError(w, http.StatusBadGateway, "upstream_error", err.Error())
			return
		}
		if !json.Valid(body) {
			msg := strings.TrimSpace(string(body))
			if msg == "" {
				msg = resp.Status
			}
			logger.Warn("upstream error body",
				"model", modelID,
				"status", resp.StatusCode,
				"body", truncate(msg, 500),
			)
			openai.WriteError(w, resp.StatusCode, "upstream_error", msg)
			return
		}
		logger.Warn("upstream json error",
			"model", modelID,
			"status", resp.StatusCode,
			"body", truncate(string(body), 500),
		)
		w.WriteHeader(resp.StatusCode)
		_, _ = w.Write(body)
		return
	}

	w.WriteHeader(resp.StatusCode)
	useStream := streaming || isEventStream(resp.Header)
	if !useStream && streaming {
		logger.Warn("client requested stream but upstream content-type is not event-stream",
			"model", modelID,
			"content_type", resp.Header.Get("Content-Type"),
		)
	}
	if flusher, ok := w.(http.Flusher); ok && useStream {
		logger.Debug("streaming response", "model", modelID)
		var flush func()
		flush = flusher.Flush
		if err := openai.RewriteSSEStream(resp.Body, w, flush); err != nil {
			if errors.Is(err, context.Canceled) || clientCtx.Err() != nil {
				logger.Debug("stream canceled by client", "model", modelID)
			} else {
				logger.Warn("stream copy failed", "model", modelID, "error", err)
			}
		}
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Warn("read upstream body failed", "model", modelID, "error", err)
		return
	}
	if !streaming {
		body, _ = openai.NormalizeChatCompletion(body)
	}
	logger.Debug("upstream body", "model", modelID, "bytes", len(body))
	_, _ = w.Write(body)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

func isEventStream(h http.Header) bool {
	return strings.Contains(strings.ToLower(h.Get("Content-Type")), "text/event-stream")
}

func retryDelay(attempt int) time.Duration {
	return time.Duration(attempt) * 500 * time.Millisecond
}
