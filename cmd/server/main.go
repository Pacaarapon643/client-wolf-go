package main

import (
	"log"
	"werewolf-backend/internal/infrastructure/config"
	"werewolf-backend/internal/infrastructure/database"
	"werewolf-backend/internal/infrastructure/server"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize database
	db, err := database.NewPostgres(cfg.DB)
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}
	defer db.Close()

	// redis
	rdb, err := database.NewRedis(cfg.Redis)
	if err != nil {
		log.Fatal("Failed to connect Redis:", err)
	}
	defer rdb.Close()
	log.Println("✅ redis connected")

	// Auto migrate
	if err := db.AutoMigrate(); err != nil {
		log.Fatal("Failed to migrate:", err)
	}

	// Initialize and run server
	srv := server.NewServer(cfg, db, rdb)
	if err := srv.Run(); err != nil {
		log.Fatal("Server error:", err)
	}
}
