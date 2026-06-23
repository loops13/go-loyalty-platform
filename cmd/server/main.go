package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"awesomeProject/internal/client"
	"awesomeProject/internal/logging"
	"awesomeProject/internal/reward"
	"awesomeProject/internal/store"
)

func main() {
	// Configuration from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	// Dependency injection: repository → services → handlers
	repo := store.NewInMemoryStore()
	clientSvc := client.NewService(repo)
	rewardSvc := reward.NewService(repo, clientSvc)

	clientHandler := client.NewHandler(clientSvc)
	rewardHandler := reward.NewHandler(rewardSvc)

	// Router setup
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(logging.Middleware(logger))
	r.Use(middleware.Recoverer)

	// Register domain handlers
	clientHandler.RegisterRoutes(r)
	rewardHandler.RegisterRoutes(r)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// HTTP server configuration
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in background
	serverErr := make(chan error, 1)
	go func() {
		logger.Info("starting server", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
			return
		}
		serverErr <- nil
	}()

	// Graceful shutdown on signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-quit:
		logger.Info("shutdown signal received", "signal", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("server forced to shutdown", "error", err)
			os.Exit(1)
		}
		logger.Info("server stopped")
	case err := <-serverErr:
		if err != nil {
			logger.Error("server stopped unexpectedly", "error", err)
			os.Exit(1)
		}
		logger.Info("server stopped")
	}
}
