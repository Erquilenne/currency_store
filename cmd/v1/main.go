package main

import (
	"currency/config"
	_ "currency/docs" // Import swagger docs
	"currency/internal/server"
	"currency/pkg/db/migrations"
	"currency/pkg/db/postgres"
	"currency/pkg/logger"
	"log"
	"os"
	"path/filepath"
)

// @title Currency API
// @version 1.0
// @description API для работы с курсами валют
// @BasePath /api/v1

func main() {
	log.Println("Starting api server")

	configPath := config.GetConfigPath(os.Getenv("config"))

	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("LoadConfig: %s", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %s", err)
	}

	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()
	appLogger.Infof("AppVersion: %s", "LogLevel: %s, Mode: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode)

	db, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		appLogger.Fatalf("Postgresql init: %s", err)
	}
	appLogger.Infof("Postgres connected, Status: %#v", db.Stats())

	migrationsPath := filepath.Join("migrations")
	if err := migrations.RunMigrations(db.DB, migrationsPath); err != nil {
		appLogger.Fatalf("Could not run migrations: %v", err)
	}
	appLogger.Info("Migrations completed successfully")

	srv := server.NewServer(cfg, db, appLogger)
	if err := srv.Run(); err != nil {
		appLogger.Fatalf("Error running server: %v", err)
	}

	defer db.Close()
}
