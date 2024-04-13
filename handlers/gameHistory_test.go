package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/joey1123455/getEasyCoins/services"
	"github.com/joey1123455/getEasyCoins/storage"
	"github.com/joey1123455/getEasyCoins/utils"
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
	cache               utils.Cache
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
	cache = *utils.NewCache()
}

// TestNewGameHistoryHandler tests the NewGameHistoryHandler function.
//
// It verifies the fields of the created instance.
func TestNewGameHistoryHandler(t *testing.T) {
	service := services.NewGameHistoryContract(client, gameHistoryContract)
	handler := NewGameHistoryHandler(service, &ctx, transactOpts, callOpts, &cache)

	// Verify the fields of the created instance
	assert.Equal(t, service, handler.services, "services field should match")
	assert.NotNil(t, ctx, handler.ctx, "ctx field should match")
	assert.Equal(t, transactOpts, handler.TransactOpts, "TransactOpts field should match")
	assert.Equal(t, callOpts, handler.CallOpts, "CallOpts field should match")
	assert.NotNil(t, handler.Cache, "Cache field should match")
}

// TestStoreGameData tests the StoreGameData function.
//
// Params:
// - t: *testing.T
func TestStoreGameData(t *testing.T) {
	service := services.NewGameHistoryContract(client, gameHistoryContract)
	handler := NewGameHistoryHandler(service, &ctx, transactOpts, callOpts, &cache)
	// Create a new HTTP request
	body := []byte(`{"gid":1234,"gtid":"test","uid":"user123","data":"some data","time":12345}`)
	req, err := http.NewRequest("POST", "/game/store", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Create a fake Gin context
	ctx, _ := gin.CreateTestContext(rr)
	ctx.Request = req

	// Call the handler function
	handler.StoreGameData(ctx)

	// Check the status code
	assert.Equal(t, http.StatusCreated, rr.Code, "expected status code 201")

	// Check the response body
	expectedBody := `{"status":"success","message":"data stored"}`
	assert.Equal(t, expectedBody, rr.Body.String(), "response body should match expected")
}

// TestGameHistory is a test function to validate the GameHistory endpoint.
//
// Params:
// - t: *testing.T
func TestGameHistory(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	service := services.NewGameHistoryContract(client, gameHistoryContract)
	handler := NewGameHistoryHandler(service, &ctx, transactOpts, callOpts, &cache)

	router := gin.New()
	router.GET("/game/history", handler.GameHistory)

	// Test case 1: Valid request
	req1, _ := http.NewRequest("GET", "/game/history?gid=123&page=1&pageSize=20", nil)
	resp1 := httptest.NewRecorder()

	router.ServeHTTP(resp1, req1)

	assert.Equal(t, http.StatusOK, resp1.Code)
	assert.Contains(t, resp1.Body.String(), "status")
	assert.Contains(t, resp1.Body.String(), "success")

	// Unmarshal the response JSON
	var jsonResponse map[string]interface{}
	err := json.Unmarshal(resp1.Body.Bytes(), &jsonResponse)
	assert.NoError(t, err)

	// Check if the 'page' key exists in the response
	assert.Contains(t, jsonResponse, "page")

	// Test case 2: Invalid page number
	req2, _ := http.NewRequest("GET", "/game/history?gid=123&page=-1&pageSize=20", nil)
	resp2 := httptest.NewRecorder()

	router.ServeHTTP(resp2, req2)

	assert.Equal(t, http.StatusBadRequest, resp2.Code)
	assert.Contains(t, resp2.Body.String(), "status")
	assert.Contains(t, resp2.Body.String(), "fail")

	// Test case 3: Invalid page size
	req3, _ := http.NewRequest("GET", "/game/history?gid=123&page=1&pageSize=-20", nil)
	resp3 := httptest.NewRecorder()

	router.ServeHTTP(resp3, req3)

	assert.Equal(t, http.StatusBadRequest, resp3.Code)
	assert.Contains(t, resp3.Body.String(), "status")
	assert.Contains(t, resp3.Body.String(), "fail")
}

// TestUserHistory tests the UserHistory handler function by sending various HTTP requests.
//
// Params:
// - t: *testing.T
// Returns nothing.
func TestUserHistory(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	service := services.NewGameHistoryContract(client, gameHistoryContract)
	handler := NewGameHistoryHandler(service, &ctx, transactOpts, callOpts, &cache)

	router := gin.New()
	router.GET("/user/history", handler.UserHistory)

	// Test case 1: Valid request
	req1, _ := http.NewRequest("GET", "/user/history?uid=uid123&page=1&pageSize=20", nil)
	resp1 := httptest.NewRecorder()

	router.ServeHTTP(resp1, req1)

	assert.Equal(t, http.StatusOK, resp1.Code)
	assert.Contains(t, resp1.Body.String(), "status")
	assert.Contains(t, resp1.Body.String(), "success")

	// Unmarshal the response JSON
	var jsonResponse map[string]interface{}
	err := json.Unmarshal(resp1.Body.Bytes(), &jsonResponse)
	assert.NoError(t, err)

	// Check if the 'page' key exists in the response
	assert.Contains(t, jsonResponse, "page")

	// Test case 2: Invalid page number
	req2, _ := http.NewRequest("GET", "/user/history?uid=uid123&page=-1&pageSize=20", nil)
	resp2 := httptest.NewRecorder()

	router.ServeHTTP(resp2, req2)

	assert.Equal(t, http.StatusBadRequest, resp2.Code)
	assert.Contains(t, resp2.Body.String(), "status")
	assert.Contains(t, resp2.Body.String(), "fail")

	// Test case 3: Invalid page size
	req3, _ := http.NewRequest("GET", "/user/history?uid=uid123&page=1&pageSize=-20", nil)
	resp3 := httptest.NewRecorder()

	router.ServeHTTP(resp3, req3)

	assert.Equal(t, http.StatusBadRequest, resp3.Code)
	assert.Contains(t, resp3.Body.String(), "status")
	assert.Contains(t, resp3.Body.String(), "fail")
}
