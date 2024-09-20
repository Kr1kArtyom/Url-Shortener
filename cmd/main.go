package main

import (
	"context"
	"flag"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"url-shortener/config"
	"url-shortener/internal/repository"
	"url-shortener/internal/repository/local"
	"url-shortener/internal/repository/postgres"
	"url-shortener/internal/server"
	"url-shortener/internal/server/handlers"
	"url-shortener/internal/services"
	"url-shortener/pkg"
)

const CONFIG_PATH = "./config/config.yaml"

func main() {
	useDatabase := flag.Bool("d", false, "use database")
	flag.Parse()

	cfg, err := config.Parse(CONFIG_PATH)
	if err != nil {
		log.Fatalf("could parse: %s", err)
	}

	log := setupLogger()

	var rep repository.Repository
	if *useDatabase {
		dbConn, err := pkg.NewDataBase(cfg.Database.DbURL)
		if err != nil {
			log.Error("could not connect to database: %s", slog.String("error", err.Error()))
			os.Exit(1)
		}
		rep = postgres.NewPgRepository(dbConn.GetDB())
		log.Info("Using Postgres repository")
	} else {
		rep = local.NewLocalRepository()
		log.Info("Using Local in-memory repository")
	}

	urlService := services.NewURLShortener(log, rep)
	handler := handlers.NewHandler(urlService, log)

	srv, err := server.NewServer(&cfg.Server, server.NewRouter(handler), log)
	if err != nil {
		log.Error("could not initialize server", slog.String("error", err.Error()))
		os.Exit(1)
	}

	sigQuit := make(chan os.Signal, 2)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)

	eg, ctx := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		s := <-sigQuit
		log.Info("captured signal", slog.String("signal", s.String()))
		return fmt.Errorf("captured signal: %v", s)
	})

	// Запуск сервера
	eg.Go(func() error {
		log.Info("Starting server")
		if err := srv.Run(ctx); err != nil {
			log.Error("error occurred while running http server", slog.String("error", err.Error()))
			return err
		}
		return nil
	})

	// Ожидание завершения (ловим ошибку либо от сигнала, либо от сервера)
	if err := eg.Wait(); err != nil {
		log.Info("gracefully shutting down the server", slog.String("reason", err.Error()))
	}

	// Остановка сервера
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("error occurred during server shutdown", slog.String("error", err.Error()))
	}
}

// убрать в другое место
func setupLogger() *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)

	return logger
}
