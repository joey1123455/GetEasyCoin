package handler

import (
	"context"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gin-gonic/gin"
	"github.com/joey1123455/getEasyCoins/data"
	"github.com/joey1123455/getEasyCoins/services"
	"github.com/joey1123455/getEasyCoins/storage"
	"github.com/joey1123455/getEasyCoins/utils"
)

type GameHistoryHandler struct {
	services     services.GameHistoryContract
	ctx          *context.Context
	TransactOpts *bind.TransactOpts
	CallOpts     *bind.CallOpts
	Cache        *utils.Cache
}

// NewGameHistoryHandler creates a new gameHistoryHandler instance.
//
// Parameters:
//
//	service: services.GameHistoryContract
//	ctx_: *context.Context
//	tans: *bind.TransactOpts
//	call: *bind.CallOpts
//	cache: *utils.Cache
//
// Return Type:
//
//	*gameHistoryHandler
func NewGameHistoryHandler(service services.GameHistoryContract, ctx_ *context.Context, tans *bind.TransactOpts, call *bind.CallOpts, cache *utils.Cache) *GameHistoryHandler {
	return &GameHistoryHandler{
		services:     service,
		ctx:          ctx_,
		TransactOpts: tans,
		CallOpts:     call,
		Cache:        cache,
	}
}

// StoreGameData godoc
// @Summary      Store game data
// @Description  stores game data in the game history handler.
// @Tags         game history
// @Accept       json
// @Produce      json
// @Param        data  body data.GameSess true  "Game data"
// @Success      200  {object}  handler.GameHistoryResOk
// @Failure      400  {object}  handler.GameHistoryResFail
// @Failure      404  {object}  handler.GameHistoryResFail
// @Failure      500  {object}  handler.GameHistoryResFail
// @Router       /game/store [post]
func (g *GameHistoryHandler) StoreGameData(ctx *gin.Context) {
	var gameSess data.GameSess
	if err := ctx.ShouldBindJSON(&gameSess); err != nil {
		response := GameHistoryResFail{
			Status:  "fail",
			Message: err.Error(),
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	tx, err := g.services.StoreGameData(g.TransactOpts, gameSess.Gid, gameSess.Gtid, gameSess.Uid, gameSess.Data, gameSess.Time)

	if err != nil {
		response := GameHistoryResFail{
			Status:  "fail",
			Message: err.Error(),
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	response := GameHistoryStoreOk{
		Status:  "success",
		Message: "transaction hex " + tx.Hash().Hex(),
	}
	ctx.JSON(http.StatusCreated, response)
}

// GameHistory godoc
// @Summary      Show game history
// @Description  handles the retrieval of game history for a given game ID. It paginates the results based on the page and pageSize query parameters.
// @Tags         game history
// @Produce      json
// @Param        gid   path      string  true  "Game ID"
// @Param        page  query     string     false  "Page number"
// @Param        pageSize  query     string     false  "Page size"
// @Success      200  {object}  handler.GameHistoryResOk
// @Failure      400  {object}  handler.GameHistoryResFail
// @Failure      404  {object}  handler.GameHistoryResFail
// @Failure      500  {object}  handler.GameHistoryResFail
// @Router       /game/history/{gid} [get]
func (g *GameHistoryHandler) GameHistory(ctx *gin.Context) {
	var res []storage.GameHistoryGameSession
	gid := ctx.Param("gid")
	println(gid)

	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		log.Println("Invalid page number")
		response := GameHistoryResFail{
			Status:  "fail",
			Message: "Invalid page number",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	pageSize, err := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))
	if err != nil || pageSize < 1 {
		log.Println("Invalid page size")
		response := GameHistoryResFail{
			Status:  "fail",
			Message: "Invalid page size",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	// gameHistory := make([]data.GameSess, 0)

	if _, found := g.Cache.Get(gid); !found {
		_gid, err := strconv.Atoi(gid)
		if err != nil {
			log.Println("while parsing gid: ", err.Error())
			response := GameHistoryResFail{
				Status:  "fail",
				Message: "internal server error",
			}
			ctx.JSON(http.StatusInternalServerError, response)
			return
		}

		res, err = g.services.GetGameData(g.CallOpts, _gid)
		if err != nil {
			log.Println("while getting game data: ", err.Error())
			response := GameHistoryResFail{
				Status:  "fail",
				Message: err.Error(),
			}
			ctx.JSON(http.StatusBadRequest, response)
			return
		}

		sort.SliceStable(res, func(i, j int) bool {
			return utils.ComparePtrFieldsDesc(&res[i], &res[j])
		})

		g.Cache.Set(gid, res, 6*time.Minute)
	}

	cachedData, _ := g.Cache.Get(gid)
	res = cachedData.([]storage.GameHistoryGameSession)
	startIndex := (page - 1) * pageSize
	endIndex := page * pageSize

	if len(res) == 0 {
		response := GameHistoryResOk{
			Status: "failed no game data for provided gid",
			Page:   []storage.GameHistoryGameSession{},
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if startIndex >= len(res) {
		response := GameHistoryResOk{
			Status: "failed no new page",
			Page:   []storage.GameHistoryGameSession{},
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	if endIndex > len(res) {
		endIndex = len(res)
	}
	response := GameHistoryResOk{
		Status: "success",
		Page:   res[startIndex:endIndex],
	}
	ctx.JSON(http.StatusOK, response)
	return
}

// UserHistory godoc
// @Summary      Show user game history
// @Description  handles the retrieval of game history for a given user ID. It paginates the results based on the page and pageSize query parameters.
// @Tags         user game history
// @Produce      json
// @Param        uid   path      string  true  "User ID"
// @Param        page  query     string     false  "Page number"
// @Param        pageSize  query     string     false  "Page size"
// @Success      200  {object}  handler.GameHistoryResOk
// @Failure      400  {object}  handler.GameHistoryResFail
// @Failure      404  {object}  handler.GameHistoryResFail
// @Failure      500  {object}  handler.GameHistoryResFail
// @Router       /game/history/user/{uid} [get]
func (g *GameHistoryHandler) UserHistory(ctx *gin.Context) {
	var res []storage.GameHistoryGameSession
	uid := ctx.Param("uid")

	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		log.Println("Invalid page number")
		response := GameHistoryResFail{
			Status:  "fail",
			Message: "Invalid page number",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	pageSize, err := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))
	if err != nil || pageSize < 1 {
		log.Println("Invalid page size")
		response := GameHistoryResFail{
			Status:  "fail",
			Message: "Invalid page size",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	// gameHistory := make([]data.GameSess, 0)

	if _, found := g.Cache.Get(uid); !found {
		if err != nil {
			log.Println("while parsing gid: ", err.Error())
			response := GameHistoryResFail{
				Status:  "fail",
				Message: "internal server error",
			}
			ctx.JSON(http.StatusInternalServerError, response)
			return
		}

		res, err = g.services.GetUserGameData(g.CallOpts, uid)
		if err != nil {
			log.Println("while getting game data: ", err.Error())
			response := GameHistoryResFail{
				Status:  "fail",
				Message: err.Error(),
			}
			ctx.JSON(http.StatusBadRequest, response)
			return
		}

		sort.SliceStable(res, func(i, j int) bool {
			return utils.ComparePtrFieldsDesc(&res[i], &res[j])
		})

		g.Cache.Set(uid, res, 6*time.Minute)
	}

	cachedData, _ := g.Cache.Get(uid)
	res = cachedData.([]storage.GameHistoryGameSession)

	if len(res) == 0 {
		response := GameHistoryResOk{
			Status: "failed no game data for provided gid",
			Page:   []storage.GameHistoryGameSession{},
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	startIndex := (page - 1) * pageSize
	endIndex := page * pageSize
	if startIndex >= len(res) {
		response := GameHistoryResOk{
			Status: "failed no new page",
			Page:   []storage.GameHistoryGameSession{},
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	if endIndex > len(res) {
		endIndex = len(res)
	}
	response := GameHistoryResOk{
		Status: "success",
		Page:   res[startIndex:endIndex],
	}
	ctx.JSON(http.StatusOK, response)
	return
}
