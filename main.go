package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"time"
)

func main() {
	//Create client in ethereum network
	//mainnet, ropsten
	client, err := ethclient.Dial("https://ropsten.infura.io/v3/2bc821ea92fd4cdeb2d18a3661e3be29")
	if err != nil {
		log.Fatal("error creating a client", err)
	}
	ctx := context.Background()

	currentBlock, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Block count: %s, %s", currentBlock.Number().String(), currentBlock.Hash().String())

	nextBlock := currentBlock.Number()
	var latestScannedBlock uint64

	ticker := time.NewTicker(2 * time.Second)
	for _ = range ticker.C {
		fmt.Printf("We are on the currentBlock {} %s\n", nextBlock.String())
		block, err := client.BlockByNumber(ctx, nextBlock)
		if err != nil {
			log.Printf("error while getting next block %s", err)
		}
		if block != nil {
			if latestScannedBlock != block.Number().Uint64() {
				fmt.Printf("Amount of transactions in a block %d\n", len(block.Transactions()))
				for _, tx := range block.Transactions() {
					fmt.Printf("TX Hash: %s\n", tx.Hash().Hex())
					latestScannedBlock = block.Number().Uint64()
				}
			}
			nextBlock = big.NewInt(int64(nextBlock.Uint64() + 1))

			log.Printf("Setting next block as {} %d", nextBlock)
		} else {
			log.Printf("No more blocks")
		}
	}
}
