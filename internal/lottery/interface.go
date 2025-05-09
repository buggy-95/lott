package lottery

// 购奖号码，通过 Bingo 判断当前号码是否是中奖号码
type BingoNum struct {
	Num   int  // 彩票号码
	Bingo bool // 是否投中
}

// 中奖结果号码，在购奖号码的基础上增加了号码类型
//
// 号码类型包括：前区胆码、前区拖码、后区胆码、后区拖码
type ResultNum struct {
	BingoNum
	Type string // 号码类型 (FrontDan: 前区胆码, FrontTuo: 前区拖码, BackDan: 后区胆码, BackTuo: 后区拖码)
}

// 购彩基本信息
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

// 复试彩票结构
type ComplexLotteryParts struct {
	LotteryBaseInfo
	FrontDan []int // 前区胆码
	FrontTuo []int // 前区拖码
	BackDan  []int // 后区胆码
	BackTuo  []int // 后区拖码
}

// 复式彩票，包含复式彩票结构和单式彩票列表
//
// 单式彩票列表由复式彩票的胆拖组合而成
type ComplexLottery struct {
	ComplexLotteryParts
	List []SingleLottery // 单式列表
}

// 单式彩票开奖结果
type SingleLotteryResult struct {
	LotteryBaseInfo
	FrontMatched int
	BackMatched  int
	Level        int
	Price        int
	Numbers      []ResultNum
}

// 复式彩票开奖结果，包含单式彩票开奖结果列表
type ComplexLotteryResult struct {
	LotteryBaseInfo
	Price   int
	Numbers []ResultNum
	List    []SingleLotteryResult
}
