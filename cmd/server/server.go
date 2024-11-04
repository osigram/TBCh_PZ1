package main

import (
	"PZ1/internal/keystorage"
	"PZ1/internal/server/handlers"
	"PZ1/internal/server/netsync"
	"PZ1/internal/server/storage/inmemory"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	// Config
	keyStorage := keystorage.MustNewKeyStorage("key.json")
	key := keyStorage.Key()

	// Logger
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(handler)

	logger.Info("Starting application...")

	// Source initialisation
	storage := inmemory.MustNewBlockchainStorage(&key)

	netsync.StartSynchronizationRoutine(logger, storage, []string{"localhost:54579"}, time.Duration(15)*time.Second)

	// Router
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/blockchain", handlers.BlockchainHandler(logger, storage))
	r.Get("/account", handlers.Account(logger, storage))
	r.Post("/transaction", handlers.Transaction(logger, storage))

	err := http.ListenAndServe(":34578", r)
	if err != nil {
		logger.Error(fmt.Sprint(err))
	}
}
