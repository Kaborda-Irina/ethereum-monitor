package services

import (
	"context"
	"crypto/ecdsa"
	"ethereum-monitor/database"
	models "ethereum-monitor/internal"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/sha3"
	"log"
	"math/big"
	"regexp"
	"strconv"
	"time"
)

// mainnet, ropsten
func ConnectionToClient() (*ethclient.Client, error) {
	client, err := ethclient.Dial("https://ropsten.infura.io/v3/2bc821ea92fd4cdeb2d18a3661e3be29")
	if err != nil {
		log.Println("error while getting connection to services ", err)
	}
	return client, err
}

func GetAddressFromDB() ([]models.DataFromDB, error) {
	db, err := database.ConnectionToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, _ := db.Query("SELECT id, address, counter FROM addresses")
	var dataFromDB models.DataFromDB
	var allDataFromDB []models.DataFromDB
	for rows.Next() {
		rows.Scan(&dataFromDB.ID, &dataFromDB.Address, &dataFromDB.Counter)
		allDataFromDB = append(allDataFromDB, dataFromDB)
	}

	return allDataFromDB, err
}

func CheckBlocks(ctx context.Context, client *ethclient.Client) {
	//address := "0x4448ebCD6Bb6DB54cce8249c6CF021EB20D822B4"
	allDataFromDB, err := GetAddressFromDB()
	if err != nil {
		log.Fatal("error getting data from db ", err)
	}

	currentBlock, err := GetBlocks(ctx, client)
	if err != nil {
		log.Fatal(err)
	}

	nextBlock := currentBlock.Number()
	var latestScannedBlock uint64
	ticker := time.NewTicker(2 * time.Second)
	go func() {
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

						msg, err := tx.AsMessage(types.LatestSignerForChainID(tx.ChainId()), nil)
						if err != nil {
							log.Printf("error while getting message %s", err)
						}
						fmt.Printf("TX from: %s\n", msg.From().Hex())
						receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
						if err != nil {
							log.Fatal("receipt: ", err)
						}
						fmt.Printf("TX status: %d\n", receipt.Status)
						fmt.Printf("TX to: %s\n", msg.To())

						for _, data := range allDataFromDB {
							go func(address string) {
								if msg.From().Hex() == address || (msg.To() != nil && msg.To().String() == address) {
									fmt.Printf("!!!!!!!!!!!!!!!!!!!!!!!!!  Address %s gets a tokens : %s", address, tx.Value().String())
									log.Fatal("completed")
								}
							}(data.Address)
						}

						latestScannedBlock = block.Number().Uint64()

					}
				}
				nextBlock = big.NewInt(int64(nextBlock.Uint64() + 1))

				log.Printf("Setting next block as {} %d", nextBlock)
			} else {
				log.Printf("No more blocks")
			}
		}
	}()
}

func GetBlocks(ctx context.Context, client *ethclient.Client) (*types.Block, error) {
	currentBlock, err := client.BlockByNumber(ctx, big.NewInt(13049312))
	if err != nil {
		log.Println("error while getting current block ", err)
	}
	log.Printf("Block count: %s, %s", currentBlock.Number().String(), currentBlock.Hash().String())

	return currentBlock, err
}

func GetBalance(ctx context.Context, client *ethclient.Client, address string) (*big.Int, error) {

	ad := common.HexToAddress(address)
	balance, err := client.BalanceAt(ctx, ad, nil)
	if err != nil {
		log.Print("There was an error", err)
		return nil, err
	}
	return balance, nil
}

func TransferringETH(client *ethclient.Client /*, privateKey *ecdsa.PrivateKey*/) {
	value := big.NewInt(100000000) // in wei (1 eth)
	gasLimit := uint64(21000)      // in units
	//gasPrice := big.NewInt(30000000000)     // in wei (30 gwei)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	toAddress := common.HexToAddress("0x8c2f16fCB4aD9072D61B543ab010462e1581E778")

	privateKey, err := crypto.HexToECDSA("4f3edf983ac636a65a842ce7c78d9aa706d3b113bce9c46f30d7d21715b23b1d")
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal("error nonce: ", err)
	}
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)
	//newTx := legacyTx()
	//tx := types.NewTx(newTx)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(" error chainID: ", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal("signedTx :", err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal("SendTransaction: ", err)
	}

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())
}
func legacyTx() *types.Transaction {
	nonce, err := hexutil.DecodeUint64("0x1216")
	if err != nil {
		panic(err)
	}
	gasPrice, err := hexutil.DecodeBig("0x2bd0875aed")
	if err != nil {
		panic(err)
	}
	gas, err := hexutil.DecodeUint64("0x5208")
	if err != nil {
		panic(err)
	}
	to := common.HexToAddress("0x2f14582947e292a2ecd20c430b46f2d27cfe213c")
	value, err := hexutil.DecodeBig("0x2386f26fc10000")
	if err != nil {
		panic(err)
	}
	data := common.Hex2Bytes("0x")
	v, err := hexutil.DecodeBig("0x1")
	if err != nil {
		panic(err)
	}
	r, err := hexutil.DecodeBig("0x56b5bf9222ce26c3239492173249696740bc7c28cd159ad083a0f4940baf6d03")
	if err != nil {
		panic(err)
	}
	s, err := hexutil.DecodeBig("0x5fcd608b3b638950d3fe007b19ca8c4ead37237eaf89a8426777a594fd245c2a")
	if err != nil {
		panic(err)
	}

	newLegacyTx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gas,
		To:       &to,
		Value:    value,
		Data:     data,
		V:        v,
		R:        r,
		S:        s,
	})
	fmt.Println("LegacyTx expected hash     => 0xb4848204c8432070136a41792003caf8dea08f9eb284eb4240845bf64a66a068")
	fmt.Println("LegacyTx actual hash       =>", newLegacyTx.Hash().String())
	return newLegacyTx
}
func GenerateAddress() string {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)
	fmt.Println(hexutil.Encode(privateKeyBytes)[2:]) //This is the private key which is used for signing transactions

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println(hexutil.Encode(publicKeyBytes)[4:]) //  This is the public key

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println(address) // address  The public address is simply the Keccak-256 hash of the public key

	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])
	return fmt.Sprintln(hexutil.Encode(hash.Sum(nil)[12:]))
}

func GenerateDeriveAddress() string {
	mnemonic := "tag volcano eight thank tide danger coast health above argue embrace heavy"
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		log.Fatal(err)
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintln(account.Address.Hex())
}

func GenerateDeriveAddress2(count int) (string, *ecdsa.PrivateKey) {
	mnemonic, err := mnemonicFun()
	if err != nil {
		log.Fatal(err)
	}

	path := "m/44'/60'/0'/0/" + strconv.Itoa(count)
	derivPath, err := accounts.ParseDerivationPath(path)
	if err != nil {
		log.Fatal(err)
	}
	// Generate a Bip32 HD wallet for the mnemonic and a user supplied password
	seed := bip39.NewSeed(*mnemonic, "")

	// Generate a new master node using the seed.
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		log.Fatal(err)
	}

	key := masterKey
	// child extended keys
	for _, n := range derivPath {
		key, err = key.Derive(n)
		if err != nil {
			log.Fatal(err)
		}
	}
	privateKey, err := key.ECPrivKey()
	privateKeyECDSA := privateKey.ToECDSA()
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal(err)
	}

	//publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	//fmt.Println(hexutil.Encode(publicKeyBytes)[4:]) //  This is the public key

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	//or
	//hash := sha3.NewLegacyKeccak256()
	//hash.Write(publicKeyBytes[1:])
	//fmt.Sprintln(hexutil.Encode(hash.Sum(nil)[12:]))

	return address, privateKeyECDSA
}

func mnemonicFun() (*string, error) {
	// Generate a mnemonic for memorization or user-friendly seeds
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return nil, err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, err
	}
	return &mnemonic, nil
}

func IsValidAddress(v string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(v)
}
