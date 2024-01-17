package main

import (
	"effectivetask/internal/config"
	"effectivetask/internal/http-server/handlers/person/identifier"
	"effectivetask/internal/http-server/middleware/logger"
	"effectivetask/internal/storage/psql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"log/slog"
	"net/http"
	"os"
)

const (
	envDebug = "debug"
	envText  = "text"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting identifier", slog.String("env", cfg.Env))

	storage, err := psql.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage: ", err)
		os.Exit(1)
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/add", identifier.New(log, storage))
	router.Post("/del", identifier.Del(log, storage))
	router.Post("/searchbyid", identifier.SearchById(log, storage))
	router.Post("/searchbyname", identifier.SearchByName(log, storage))
	router.Post("/editbyid", identifier.Edit(log, storage))

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envDebug:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envText:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
