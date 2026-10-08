package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/masseyis/interviewtopia-completion-service/internal/appconfig"
	"github.com/masseyis/interviewtopia-completion-service/internal/httpapi"
	"github.com/masseyis/interviewtopia-completion-service/pkg/evidence"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	config := appconfig.FromEnv()
	trustStoreJSON, err := os.ReadFile(config.TrustStorePath)
	if err != nil {
		slog.Error("read trust store", "error", err)
		os.Exit(1)
	}
	trustStore, err := evidence.LoadTrustStore(trustStoreJSON)
	if err != nil {
		slog.Error("load trust store", "error", err)
		os.Exit(1)
	}

	server := &http.Server{
		Addr: config.HTTPAddr,
		Handler: httpapi.NewHandler(httpapi.Dependencies{
			RegistryURL: config.RegistryURL,
			TrustStore:  trustStore,
			HTTPClient:  &http.Client{Timeout: 3 * time.Second},
		}),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		slog.Info("completion service listening", "address", config.HTTPAddr)
		errCh <- server.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server stopped unexpectedly", "error", err)
			os.Exit(1)
		}
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			slog.Error("graceful shutdown failed", "error", err)
			os.Exit(1)
		}
	}
}
