package main

import (
	"log"

	"github.com/shooooooma415/guess-title-game-api/config"
	"github.com/shooooooma415/guess-title-game-api/internal/interface/handler"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database (commented out until repositories are implemented)
	// dbCfg := persistence.Config{
	// 	Host:     cfg.Database.Host,
	// 	Port:     cfg.Database.Port,
	// 	User:     cfg.Database.User,
	// 	Password: cfg.Database.Password,
	// 	DBName:   cfg.Database.DBName,
	// 	SSLMode:  cfg.Database.SSLMode,
	// }
	// db, err := persistence.NewDB(dbCfg)
	// if err != nil {
	// 	log.Fatalf("Failed to connect to database: %v", err)
	// }
	// defer db.Close()

	// Initialize router
	e := handler.NewRouter()

	// Start server
	log.Printf("Starting server on port %s", cfg.Server.Port)
	if err := e.Start(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
