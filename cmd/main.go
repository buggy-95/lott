package main

import (
	"github.com/buggy-95/lott/internal/lottery"
)

func main() {
	targetLottery := "DLT:15,18,20,21,34-04,10"
	sourceLotteryList := []string{
		"DLT:11,14,15,31,35-01,11x3",
		"DLT:02,08,09,15,30-05,08x3",
		"DLT:13,15,16,20,22-10,12x3",
		"DLT:09,15,22,23,31-07,08x3",
		"DLT:05,16,19,31,34-04,10x3",
	}

	target, _ := lottery.GetLottery(targetLottery)

	for _, str := range sourceLotteryList {
		lott, _ := lottery.GetLottery(str)
		result, _ := lott.GetLotteryResult(target)

		result.PrintResult(true, true)
		result.PrintList(true, true)
	}
}
