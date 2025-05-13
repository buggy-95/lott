package dlt

type PrizeLevel struct {
	AwardType         int    `json:"awardType"`
	Group             string `json:"group"`
	PrizeLevel        string `json:"prizeLevel"`
	Sort              int    `json:"sort"`
	StakeAmount       string `json:"stakeAmount"`
	StakeAmountFormat string `json:"stakeAmountFormat"`
	StakeCount        string `json:"stakeCount"`
	TotalPrizeamount  string `json:"totalPrizeamount"`
}

type PoolDraw struct {
	LotteryDrawNum       string       `json:"lotteryDrawNum"`
	LotteryDrawResult    string       `json:"lotteryDrawResult"`
	LotteryDrawTime      string       `json:"lotteryDrawTime"`
	PoolBalanceAfterdraw string       `json:"poolBalanceAfterdraw"`
	PrizeLevelList       []PrizeLevel `json:"prizeLevelList"`
}

type HistoryValue struct {
	LastPoolDraw PoolDraw   `json:"lastPoolDraw"`
	List         []PoolDraw `json:"list"`
	PageNo       int        `json:"pageNo"`
	PageSize     int        `json:"pageSize"`
	Pages        int        `json:"pages"`
	Total        int        `json:"total"`
}

type HistoryResponse struct {
	ErrorCode    string       `json:"errorCode"`
	ErrorMessage string       `json:"errorMessage"`
	Success      bool         `json:"success"`
	Value        HistoryValue `json:"value"`
}
