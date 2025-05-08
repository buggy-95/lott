package lottery

type BingoNum struct {
	Num   int  // 彩票号码
	Bingo bool // 是否投中
}

type ResultNum struct {
	BingoNum
	Type string // 号码类型 (FrontDan: 前区胆码, FrontTuo: 前区拖码, BackDan: 后区胆码, BackTuo: 后区拖码)
}

type LotteryResult struct {
	Source SingleLottery // 兑奖彩票
	Target SingleLottery // 开奖彩票
	Level  int           // 中奖等级
	Price  int           // 中奖金额
}

type LotteryBaseInfo struct {
	Type  string // 彩票类型 (DLT: 大乐透, SSQ: 双色球)
	Index int    // 开奖期号
	Scale int    // 倍投倍数
}

// 单式彩票
type SingleLottery struct {
	LotteryBaseInfo
	Front []int // 前区号码
	Back  []int // 后区号码
}

// 复试彩票
type ComplexLotteryParts struct {
	LotteryBaseInfo
	FrontDan []int // 前区胆码
	FrontTuo []int // 前区拖码
	BackDan  []int // 后区胆码
	BackTuo  []int // 后区拖码
}

type ComplexLottery struct {
	ComplexLotteryParts
	List []SingleLottery // 单式列表
}

type SingleLotteryResult struct {
	LotteryBaseInfo
	Numbers      []ResultNum
	FrontMatched int
	BackMatched  int
	Level        int
	Price        int
}

type ComplexLotteryResult struct {
	LotteryBaseInfo
	Price   int
	Numbers []ResultNum
	List    []SingleLotteryResult
}
