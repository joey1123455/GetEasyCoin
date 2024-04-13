package main

import (
	"context"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joey1123455/getEasyCoins/config"
	docs "github.com/joey1123455/getEasyCoins/docs"
	handler "github.com/joey1123455/getEasyCoins/handlers"
	"github.com/joey1123455/getEasyCoins/middleware"
	"github.com/joey1123455/getEasyCoins/routes"
	"github.com/joey1123455/getEasyCoins/services"
	"github.com/joey1123455/getEasyCoins/storage"
	"github.com/joey1123455/getEasyCoins/utils"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	ctx                 context.Context
	client              *ethclient.Client
	contractAddress     common.Address
	gameHistoryContract *storage.GameHistory
	transactOpts        *bind.TransactOpts
	callOpts            *bind.CallOpts
	gameHistoryService  services.GameHistoryContract
	gameHistoryHandler  handler.GameHistoryHandler
	gameHistoryRouter   routes.GameDataRouteController
	server              *gin.Engine
	cache               utils.Cache
)

func main() {
	config, err := config.LoadConfig(".")

	if err != nil {
		panic("Could not load config: " + err.Error())
	}
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{config.ORIGIN}
	corsConfig.AllowCredentials = true
	server.Use(cors.New(corsConfig))
	server.Use(middleware.RecoveryWithFileLogger("logs/panic.log"))

	docs.SwaggerInfo.Title = "Easy Get Coin Leader Board API"
	docs.SwaggerInfo.Description = "The leader board for easy get coin games."
	docs.SwaggerInfo.Version = config.API_VERSION
	docs.SwaggerInfo.Host = config.API_HOST + ":" + config.PORT
	docs.SwaggerInfo.BasePath = "/api"

	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "ok"})
	})
	gameHistoryRouter.GameDataRoute(router)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	log.Fatal(server.Run(":" + config.PORT))
}

func init() {
	config, err := config.LoadConfig(".")
	if err != nil {
		panic("COULD NOT LOAD CONFIG: " + err.Error())
	}

	ctx = context.TODO()
	client, err = ethclient.Dial(config.NODE_URL)
	if err != nil {
		panic("Failed to connect to the Ethereum client: " + err.Error())
	}

	contractAddress = common.HexToAddress(config.CONTRACT_ADDRESS)
	// gameHistoryContract, err = storage.NewStorage(contractAddress, client)
	// contractCaller = storage.NewGameHistoryCaller()
	privateKey, err := crypto.HexToECDSA(config.PRIVATE_KEY)
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	transactOpts = bind.NewKeyedTransactor(privateKey)

	gameHistoryContract, err = storage.NewGameHistory(contractAddress, client)
	if err != nil {
		panic("Failed to instantiate contract: " + err.Error())
	}

	callOpts = &bind.CallOpts{Context: ctx}
	cache = *utils.NewCache()

	gameHistoryService = services.NewGameHistoryContract(client, gameHistoryContract)
	gameHistoryHandler = *handler.NewGameHistoryHandler(gameHistoryService, &ctx, transactOpts, callOpts, &cache)
	gameHistoryRouter = routes.NewGameDataRouteController(gameHistoryHandler)
	server = gin.Default()
	gin.SetMode(config.MODE)
}
