package main

import (
	"context"
	"os"

	kHTTP "github.com/ObscuraNote/api-general/internal/keys/http"
	keysRepository "github.com/ObscuraNote/api-general/internal/keys/repository"
	keysService "github.com/ObscuraNote/api-general/internal/keys/service"
	uHTTP "github.com/ObscuraNote/api-general/internal/users/http"
	usersRepository "github.com/ObscuraNote/api-general/internal/users/repository"
	userService "github.com/ObscuraNote/api-general/internal/users/service"
	"github.com/ObscuraNote/api-general/internal/utils/config"
	"github.com/philippe-berto/database/postgresdb"
	httpkit "github.com/philippe-berto/httpkit"
	metrics "github.com/philippe-berto/httpkit/metrics"
	"github.com/philippe-berto/logger"
	"github.com/philippe-berto/tracer"
)

var migrationsPath = "file://./migrations"

func main() {
	ctx := context.Background()
	log := logger.New(ctx)
	cfg, err := config.Load()
	if err != nil {
		log.WithFields(logger.Fields{"error": err.Error(), "component": "main", "function": "main"}).
			Error("Failed to load configuration")

		os.Exit(1)
	}
	log.Info("Load configuration: %+v", cfg)

	log.SetLevel(cfg.Debug)

	closeTracer, err := tracer.New(ctx, cfg.Tracer, cfg.Service, cfg.Name)
	if err != nil {
		os.Exit(1)
	}

	db, err := postgresdb.New(ctx, cfg.DB, false, migrationsPath)
	if err != nil {
		log.WithFields(logger.Fields{"error": err.Error(), "component": "main", "function": "main"}).
			Error("Failed to connect to database")

		os.Exit(1)
	}

	uRepo, err := usersRepository.New(ctx, db)
	if err != nil {
		log.WithFields(logger.Fields{"error": err.Error(), "component": "main", "function": "main"}).
			Error("Failed to create users repository")
		os.Exit(1)
	}
	log.Info("Users repository initialized")

	uServ := userService.New(ctx, uRepo)
	log.Info("User service initialized")

	kRepo := keysRepository.New(ctx, db)
	if kRepo == nil {
		log.WithFields(logger.Fields{"error": "Failed to create keys repository", "component": "main", "function": "main"}).
			Error("Failed to create keys repository")
		os.Exit(1)
	}

	kServ := keysService.New(ctx, *log, kRepo, uServ)
	log.Info("Keys service initialized")

	server := httpkit.New(cfg.Port, false, false, cfg.EnableCORS, cfg.CorsAllowOrigins)
	uHTTP.Register(server.Router, uServ, *log)
	kHTTP.Register(server.Router, &kServ, uServ, *log)

	go metrics.StartMetrics(cfg.Metrics.Port, cfg.Metrics.Enable, log)

	log.Info("Starting HTTP Server at port: %d", cfg.Port)
	if err := server.Start(); err != nil {
		log.WithFields(logger.Fields{"error": err.Error()}).
			Error("HTTP Server: error on starting HTTP Server")

		os.Exit(1)
	}

	if err := server.GracefulShutdown(ctx, 60); err != nil {
		log.WithFields(logger.Fields{"error": err.Error()}).
			Error("HTTP Server: Error on shutting down HTTP Gracefully")

		os.Exit(1)
	}

	log.Info("Closing HTTP Server")
	if err := closeTracer(ctx); err != nil {
		log.WithFields(logger.Fields{"error": err.Error()}).
			Error("Tracer: Error closing down")

		os.Exit(1)
	}

}
