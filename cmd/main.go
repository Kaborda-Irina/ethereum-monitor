package main

import (
	"context"
	"ethereum-monitor/database"
	"ethereum-monitor/internal/handlers"
	"ethereum-monitor/internal/server"
	"ethereum-monitor/internal/services"
	"ethereum-monitor/vault"
	"log"
)

func main() {
	//Create client in ethereum services
	client, err := services.ConnectionToClient()
	if err != nil {
		log.Fatal("error creating a client", err)
	}

	// Start Vault
	dataVault := vault.InitVault()

	ctx := context.Background()
	accountIndex := counter()
	dataHandler := handlers.NewHandler(ctx, client, accountIndex, dataVault)

	// Check wallets in network
	services.CheckBlocks(ctx, client)

	// Start server
	server.InitServer(dataHandler)

}

func counter() int {
	db, err := database.ConnectionToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var count int
	err = db.QueryRow("SELECT max(counter) FROM addresses").Scan(&count)
	if err != nil {
		count = 0
		log.Println(err)
	}
	return count
}
