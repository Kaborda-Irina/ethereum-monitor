package handlers

import (
	"context"
	"encoding/json"
	"ethereum-monitor/database"
	models "ethereum-monitor/internal"
	"ethereum-monitor/internal/services"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Data struct {
	Ctx    context.Context
	Client *ethclient.Client
	Count  int
}

func NewHandler(ctx context.Context, client *ethclient.Client, count int) *Data {
	return &Data{
		Ctx:    ctx,
		Client: client,
		Count:  count,
	}
}

func (d *Data) AddAddress(w http.ResponseWriter, _ *http.Request) {
	address, privateKey := services.GenerateDeriveAddress2(d.Count)

	privateKeys := string(crypto.FromECDSA(privateKey))
	d.Count++

	db, err := database.ConnectionToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var counter int
	counter++
	statement, _ := db.Prepare("INSERT INTO addresses (address,privateKeys, counter) VALUES (?, ?, ?)")
	statement.Exec(address, privateKeys, counter)
	rows, _ := db.Query("SELECT id, address, privateKeys, counter FROM addresses")
	var dataFromDB models.DataFromDB
	var allDataFromDB []models.DataFromDB
	for rows.Next() {
		rows.Scan(&dataFromDB.ID, &dataFromDB.Address, &dataFromDB.PrivateKeys, &dataFromDB.Counter)
		allDataFromDB = append(allDataFromDB, dataFromDB)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&allDataFromDB)
	if err != nil {
		log.Println(err)
	}

}

func (d *Data) GetAddress(w http.ResponseWriter, _ *http.Request) {
	allDataFromDB, err := services.GetAddressFromDB()
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&allDataFromDB)
	if err != nil {
		log.Println(err)
	}
}

func (d *Data) GetBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	balance, err := services.GetBalance(d.Ctx, d.Client, address)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&balance)
	if err != nil {
		log.Println(err)
	}
}
