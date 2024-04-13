package handler

type GameHistoryResOk struct {
	Status string `json:"status"`
	Page   any    `json:"page"`
}

type GameHistoryStoreOk struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type GameHistoryResFail struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
