package main

import (
	"context"
	"hackaton/internal/config"
	"hackaton/internal/handler"
	"hackaton/internal/service"
	"hackaton/internal/storage/mongodb"
	"hackaton/internal/transport/router"
	"hackaton/internal/utils/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config := config.NewConfig()

	appLogger := logger.SetupLogger(config.Env)

	appLogger.Info().Msg("Application started")
	appLogger.Info().Str("env", config.Env).Msg("Environment loaded")

	repo, err := mongodb.NewStorage(
		config.DBConnectionString,
		"rostov_travel",
		"users",
	)

	if err != nil {
		appLogger.Fatal().Err(err).Msg("Failed to connect to database")
	}

	appLogger.Info().Msg("Database connect success")

	service := service.NewService(repo, config.JWTSecret, appLogger)

	appLogger.Info().Msg("Service init success")

	handler := handler.NewHandler(service)

	appLogger.Info().Msg("Handler init success")

	server := router.NewServer([]byte(config.JWTSecret), handler, config)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Error().Err(err).Msg("Server error")
		}
	}()

	appLogger.Info().Str("port", config.HTTPServer.Port).Msg("Server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		appLogger.Error().Err(err).Msg("Server shutdown error")
	}

	appLogger.Info().Msg("Server exited")
}
