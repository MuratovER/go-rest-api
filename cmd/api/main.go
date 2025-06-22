package main

import (
	"inventory_assets/config"
	"inventory_assets/internal/server"
	"inventory_assets/pkg/db/postgres"
	"inventory_assets/pkg/logger"
	"inventory_assets/pkg/utils"
	"log"
	"os"
)

// @title Inventory Assets API
// @version 1.0
// @description REST API for managing inventory assets
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /api/v1
func main() {
	log.Println("Starting Inventory Assets API server")

	// Load configuration
	configPath := utils.GetConfigPath(os.Getenv("config"))
	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	// Initialize logger
	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s",
		cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode)

	// Initialize database connection
	psqlDB, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		appLogger.Fatalf("Failed to initialize PostgreSQL: %s", err)
	}

	db, err := psqlDB.DB()
	if err != nil {
		appLogger.Fatalf("Failed to get database instance: %s", err)
	}
	defer db.Close()

	appLogger.Infof("PostgreSQL connected, Status: %#v", db.Stats())

	// Initialize and start server
	server := server.NewServer(cfg, psqlDB, appLogger)
	if err = server.Run(); err != nil {
		appLogger.Fatalf("Failed to start server: %v", err)
	}
}
