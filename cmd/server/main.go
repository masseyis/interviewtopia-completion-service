package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/masseyis/interviewtopia-completion-service/internal/appconfig"
	"github.com/masseyis/interviewtopia-completion-service/internal/httpapi"
	"github.com/masseyis/interviewtopia-completion-service/internal/storage"
	"github.com/masseyis/interviewtopia-completion-service/pkg/evidence"
)

func main() {
	if err := run(); err != nil {
		slog.Error("completion service stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	config := appconfig.FromEnv()
	trustStoreJSON, err := os.ReadFile(config.TrustStorePath)
	if err != nil {
		return fmt.Errorf("read trust store: %w", err)
	}
	trustStore, err := evidence.LoadTrustStore(trustStoreJSON)
	if err != nil {
		return fmt.Errorf("load trust store: %w", err)
	}
	database, err := storage.Open(config.DatabasePath)
	if err != nil {
		return err
	}
	defer database.Close()

	server := &http.Server{
		Addr: config.HTTPAddr,
		Handler: httpapi.NewHandler(httpapi.Dependencies{
			RegistryBaseURL: config.RegistryBaseURL,
			TrustStore:      trustStore,
			Database:        database,
			HTTPClient:      &http.Client{Timeout: 3 * time.Second},
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
			return fmt.Errorf("listen: %w", err)
		}
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("graceful shutdown: %w", err)
		}
	}
	return nil
}
