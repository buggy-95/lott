package lottery

type Lottery struct {
	Type  string // 彩票类型（DLT: 大乐透，SSQ: 双色球）
	Index int    // 开奖期号
	Scale int    // 倍投倍数
	Red   []int  // 红球
	Blue  []int  //蓝球
}

type LotteryResult struct {
	Source Lottery // 兑奖彩票
	Target Lottery // 开奖彩票
	Level  int     // 中奖等级
	Price  int     // 中奖金额
}
