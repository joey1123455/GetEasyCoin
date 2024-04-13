package services

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joey1123455/getEasyCoins/storage"
)

type GameHistoryContract interface {
	StoreGameData(transactData *bind.TransactOpts, gid int, gtid string, uid string, data string, time int) (res *types.Transaction, err error)
	GetGameData(callData *bind.CallOpts, gid int) (res []storage.GameHistoryGameSession, err error)
	GetUserGameData(callData *bind.CallOpts, uid string) (res []storage.GameHistoryGameSession, err error)
}

type gameHistory struct {
	ethClient *ethclient.Client
	contract  *storage.GameHistory
}

// NewGameHistoryContract creates a new instance of GameHistoryContract.
//
// Params:
// - client: the Ethereum client
// - contract: the storage contract
// - tans: the transaction options
// - call: the call options
// Returns a GameHistoryContract instance.
func NewGameHistoryContract(client *ethclient.Client, contract *storage.GameHistory) GameHistoryContract {
	return &gameHistory{
		ethClient: client,
		contract:  contract,
	}
}

// StoreGameData stores game data in the game history.
//
// Params:
// - transactData: transaction options for the contract.
// - gid: game ID.
// - gtid: game transaction ID.
// - uid: user ID.
// - data: game data.
// - time: timestamp.
//
// Returns:
// - res: transaction result.
// - err: error.
func (g *gameHistory) StoreGameData(transactData *bind.TransactOpts, gid int, gtid string, uid string, data string, time int) (res *types.Transaction, err error) {
	ctx := context.TODO()
	suggestedFee, err := g.ethClient.SuggestGasPrice(ctx)
	if err != nil {
		suggestedFee = big.NewInt(200000000)
	}

	suggestedTip, err := g.ethClient.SuggestGasTipCap(ctx)
	if err != nil {
		suggestedTip = big.NewInt(200000)
	}

	transactData.Value = suggestedTip
	transactData.GasLimit = uint64(3000000)
	transactData.GasPrice = suggestedFee
	log.Println(transactData.GasPrice)
	log.Println(transactData.GasLimit)
	// res, err = g.contract.StoreGameData(transactData, big.NewInt(int64(gid)), gtid, uid, data, big.NewInt(int64(time)))
	res, err = g.contract.StoreGameData(transactData, big.NewInt(int64(gid)), gtid, uid, data, big.NewInt(int64(time)))
	return
}

// GetGameData retrieves game data for a specific game ID.
//
// Params:
// - callData - Call options for the contract call.
// - gid - The game ID for which data is requested.
//
// Returns:
// - res : an array of game history game sessions or an error.
// - error : Any error that occurred during the retrieval.
func (g *gameHistory) GetGameData(callData *bind.CallOpts, gid int) (res []storage.GameHistoryGameSession, err error) {
	res, err = g.contract.GetGameHistory(callData, big.NewInt(int64(gid)))
	return
}

// GetUserGameData retrieves game data for a specific user.
//
// Parameters:
//
// - callData *bind.CallOpts: Options for the contract call
// - uid string: User ID for which to retrieve the game data
//
// Returns:
//
// - res []storage.GameHistoryGameSession: Game data for the user
// - err error: Any error that occurred during the retrieval
func (g *gameHistory) GetUserGameData(callData *bind.CallOpts, uid string) (res []storage.GameHistoryGameSession, err error) {
	res, err = g.contract.GetUserHistory(callData, uid)
	return
}
