package services

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joey1123455/getEasyCoins/storage"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var (
	ctx                 context.Context
	client              *ethclient.Client
	contractAddress     common.Address
	gameHistoryContract *storage.GameHistory
	transactOpts        *bind.TransactOpts
	callOpts            *bind.CallOpts
)

func init() {
	_ = godotenv.Load("../.env")
	nodeUrl := os.Getenv("NODE_URL")
	address := os.Getenv("CONTRACT_ADDRESS")
	privateKeyValue := os.Getenv("PRIVATE_KEY")

	ctx = context.TODO()
	client, err := ethclient.Dial(nodeUrl)
	if err != nil {
		panic("Failed to connect to the Ethereum client: " + err.Error())
	}

	contractAddress = common.HexToAddress(address)
	gameHistoryContract, err = storage.NewGameHistory(contractAddress, client)
	if err != nil {
		panic("Failed to instantiate contract: " + err.Error())
	}

	privateKey, err := crypto.HexToECDSA(privateKeyValue) // Replace with your actual private key
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	transactOpts = bind.NewKeyedTransactor(privateKey)
	callOpts = &bind.CallOpts{Context: ctx}
}

func TestNewGameHistoryContract(t *testing.T) {
	gameHistoryContract := NewGameHistoryContract(client, gameHistoryContract)
	_, ok := gameHistoryContract.(*gameHistory)
	assert.True(t, ok, "expected gameHistory type")
}

func TestStoreGameData(t *testing.T) {
	g := NewGameHistoryContract(client, gameHistoryContract)
	res, err := g.StoreGameData(transactOpts, 123456, "gtid123", "uid123", "some_data", 123456)
	assert.NoError(t, err, "unexpected error")
	assert.NotNil(t, res.Hash().Hex(), "transaction result should not be nil")
}

func TestGetUserGameData(t *testing.T) {
	g := NewGameHistoryContract(client, gameHistoryContract)
	res, err := g.GetUserGameData(callOpts, "uid123")
	// fmt.Print(address)
	// Assertions
	assert.NoError(t, err, "unexpected error")
	assert.NotNil(t, res, "transaction result should not be nil")
	t.Log(err)
}

func TestGetGameData(t *testing.T) {
	g := NewGameHistoryContract(client, gameHistoryContract)
	res, err := g.GetGameData(callOpts, 123)
	// fmt.Print(address)
	// Assertions
	assert.NoError(t, err, "unexpected error")
	assert.NotNil(t, res, "transaction result should not be nil")
	t.Log(err)
}
