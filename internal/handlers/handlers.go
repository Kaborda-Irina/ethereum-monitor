package handlers

import (
	"context"
	"encoding/json"
	"ethereum-monitor/database"
	models "ethereum-monitor/internal"
	"ethereum-monitor/internal/services"
	"ethereum-monitor/vault"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type Data struct {
	Ctx          context.Context
	Client       *ethclient.Client
	AccountIndex int
	DataVault    *vault.DataVault
	SecretData   map[string]int
}

func NewHandler(ctx context.Context, client *ethclient.Client, accountIndex int, dataVault *vault.DataVault) *Data {
	return &Data{
		Ctx:          ctx,
		Client:       client,
		AccountIndex: accountIndex,
		DataVault:    dataVault,
	}
}

func (d *Data) AddAddress(w http.ResponseWriter, _ *http.Request) {
	d.AccountIndex++
	address, privateKey := services.GenerateDeriveAddress(d.AccountIndex)
	services.IsValidAddress(address)

	db, err := database.ConnectionToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	counter := d.AccountIndex

	statement, _ := db.Prepare("INSERT INTO addresses (address, counter) VALUES (?, ?)")
	statement.Exec(address, counter)
	rows, _ := db.Query("SELECT id, address, counter FROM addresses")
	var dataFromDB models.DataFromDB
	var allDataFromDB []models.DataFromDB
	for rows.Next() {
		rows.Scan(&dataFromDB.ID, &dataFromDB.Address, &dataFromDB.Counter)
		allDataFromDB = append(allDataFromDB, dataFromDB)
	}

	secretData := map[string]interface{}{
		strconv.Itoa(counter): privateKey,
	}
	err = vault.WriteKey(d.DataVault, secretData)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, "Secret %s written successfully \nAddresses: \n", strconv.Itoa(counter))
	for _, data := range allDataFromDB {
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(&data)
		if err != nil {
			log.Println(err)
		}
	}
}

func (d *Data) TransferETH(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountIndex := vars["accountIndex"]
	privateKey := vault.ReadKey(d.DataVault, accountIndex)
	txHash := services.TransferringETH(d.Client, privateKey)

	fmt.Fprintf(w, "Secret read successfully \nTransaction send %s\n", txHash)
}

func (d *Data) GetAddress(w http.ResponseWriter, _ *http.Request) {
	allDataFromDB, err := services.GetAddressFromDB()
	if err != nil {
		log.Println(err)
	}

	fmt.Fprintf(w, "Getting all addresses:\n")
	for _, data := range allDataFromDB {
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(&data)
		if err != nil {
			log.Println(err)
		}
	}
}

func (d *Data) GetBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	balance, err := services.GetBalance(d.Ctx, d.Client, address)
	if err != nil {
		log.Println(err)
	}
	eth := services.WeiToEther(balance)
	fmt.Println("eth", eth)
	fmt.Fprintf(w, "Balance %s wei = %s  or eth = %f\n", address, balance, eth)
}
