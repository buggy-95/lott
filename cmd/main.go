package main

import (
	"fmt"

	"github.com/buggy-95/lott/internal/lottery"
)

func main() {
	// target := "15,18,20,21,34-04,10"
	// sources := []string{
	// 	"11,14,15,31,35-01,11",
	// 	"02,08,09,15,30-05,08",
	// 	"13,15,16,20,22-10,12",
	// 	"09,15,22,23,31-07,08",
	// 	"05,16,19,31,34-04,10",
	// }
	// resultMap := map[int]string{
	// 	1: "一等奖",
	// 	2: "二等奖",
	// 	3: "三等奖",
	// 	4: "四等奖",
	// 	5: "五等奖",
	// 	6: "六等奖",
	// 	7: "七等奖",
	// 	8: "八等奖",
	// 	9: "九等奖",
	// }

	target, _ := lottery.GetLottery("DLT:04,12,17,23,27-02,08")
	lott, _ := lottery.GetLottery("DLT:09,22,04~07,11,13,16,26,27,32,24-04~07,09,10")
	result, _ := lott.GetLotteryResult(target)

	fmt.Printf("level: %d, size: %d, price: %d", result.Level, len(result.List), result.Price)
}
