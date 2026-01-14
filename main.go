package main

import (
	"werewolf-backend/internal/config"
	"werewolf-backend/internal/server"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	server := server.New(cfg)
	if err := server.Run(); err != nil {
		panic(err)
	}
	

}
