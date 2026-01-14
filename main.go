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

	server := server.New(cfg, db)
	if err := server.Run(); err != nil {
		panic(err)
	}
	

}
