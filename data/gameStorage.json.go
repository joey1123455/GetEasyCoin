package data

type GameSess struct {
	Gid  int    `json:"gid" binding:"required"`
	Gtid string `json:"gtid" binding:"required"`
	Uid  string `json:"uid" binding:"required"`
	Data string `json:"data" binding:"required"`
	Time int    `json:"time" binding:"required"`
}
