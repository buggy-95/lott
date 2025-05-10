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

// 彩票构成部分
type LotteryParts struct {
	LotteryBaseInfo
	FrontDan []int // 前区胆码
	FrontTuo []int // 前区拖码
	BackDan  []int // 后区胆码
	BackTuo  []int // 后区拖码
}

// 彩票结构，包含组成部分和列表，若列表为空则当前彩票为单式票，复式票的列表会包含所有组成的单式票
type Lottery struct {
	LotteryParts
	List []Lottery // 单式列表，若为空则当前彩票为单式，否则为复式
}

// 彩票开奖结果
//
// 若为单式票，列表为空
//
// 若为复试票，Level为单式票列表中最高中奖等级
type LotteryResult struct {
	LotteryBaseInfo // TODO: 改成指针
	FrontMatched    int
	BackMatched     int
	Level           int
	Price           int
	Numbers         []ResultNum
	List            []LotteryResult
}
