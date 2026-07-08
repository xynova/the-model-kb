package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/xynova/cloudrun-proxy/internal/config"
	"github.com/xynova/cloudrun-proxy/internal/logging"
	"github.com/xynova/cloudrun-proxy/internal/openai"
	"github.com/xynova/cloudrun-proxy/internal/upstream"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	logger := logging.New(cfg.Server.Debug)

	ctx := context.Background()
	router, err := upstream.NewRouter(ctx, cfg, logger)
	if err != nil {
		logger.Error("router", "error", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("GET /v1/models", func(w http.ResponseWriter, r *http.Request) {
		if err := openai.CheckLocalSecret(r, cfg.Server.LocalSecret); err != nil {
			openai.WriteError(w, http.StatusUnauthorized, "authentication_error", err.Error())
			return
		}
		openai.WriteModels(w, cfg.ModelIDs())
	})
	mux.HandleFunc("POST /v1/chat/completions", router.ChatCompletions)

	handler := logging.AccessLog(logger, mux)

	srv := &http.Server{
		Addr:              cfg.Server.ListenAddr,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		logArgs := []any{
			"addr", cfg.Server.ListenAddr,
			"models", strings.Join(cfg.ModelIDs(), ", "),
			"debug", cfg.Server.Debug,
		}
		if cfg.ConfigFile != "" {
			logArgs = append(logArgs, "config", cfg.ConfigFile)
		}
		logger.Info("listening", logArgs...)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server", "error", err)
			os.Exit(1)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}
