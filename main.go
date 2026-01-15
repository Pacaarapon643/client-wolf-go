package main

import (
	"log"
	"werewolf-backend/internal/config"
	"werewolf-backend/internal/database"
	"werewolf-backend/internal/server"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	db, err := database.NewPostgres(cfg.DB)
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}
	defer db.Close()

	// Auto migrate
	if err := db.AutoMigrate(); err != nil {
		log.Fatal("Failed to migrate:", err)
	}

	server := server.NewServer(cfg, db)
	if err := server.Run(); err != nil {
		panic(err)
	}

}
