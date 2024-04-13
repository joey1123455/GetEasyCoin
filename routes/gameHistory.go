package routes

import (
	"github.com/gin-gonic/gin"
	handler "github.com/joey1123455/getEasyCoins/handlers"
)

type GameDataRouteController struct {
	gameHistoryHandler handler.GameHistoryHandler
}

func NewGameDataRouteController(gameHistoryHandler handler.GameHistoryHandler) GameDataRouteController {
	return GameDataRouteController{gameHistoryHandler}
}

// GameDataRoute handles the routes related to game data.
//
// Takes in a gin.RouterGroup as a parameter and does not return anything.
func (r *GameDataRouteController) GameDataRoute(rg *gin.RouterGroup) {
	router := rg.Group("/game")

	router.POST("/store", r.gameHistoryHandler.StoreGameData)
	router.GET("/history/:gid", r.gameHistoryHandler.GameHistory)
	router.GET("/history/user/:uid", r.gameHistoryHandler.UserHistory)
}
