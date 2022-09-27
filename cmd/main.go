package main

import (
	"context"
	"ethereum-monitor/internal/handlers"
	"ethereum-monitor/internal/server"
	"ethereum-monitor/internal/services"
	"log"
)

func main() {
	//Create client in ethereum services
	client, err := services.ConnectionToClient()
	if err != nil {
		log.Fatal("error creating a client", err)
	}
	var accountIndex int
	ctx := context.Background()
	dataHandler := handlers.NewHandler(ctx, client, accountIndex)

	// Check wallets in network
	services.CheckBlocks(ctx, client)

	// Start server
	server.InitServer(dataHandler)

	//address, countfrom := services.GenerateDeriveAddress2(accountIndex)
	//fmt.Printf("new address %s coun %d\n", address, countfrom)

	//nonce, countfrom, privateKey := services.GenerateDeriveAddress3(client, count)
	//count = countfrom
	//fmt.Println("nonce ", nonce, " count ", count, "privateKey ", privateKey)
	//ctx := context.Background()
	//services.GetBalance(client, ctx)
	//handlers.AddCoins()

	//services.TransferringETH(client)
}
