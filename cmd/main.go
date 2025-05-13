package main

import (
	"github.com/buggy-95/lott/internal/lottery"
	"github.com/buggy-95/lott/internal/lottery/dlt"
)

func main() {
	targetLottery := "DLT:02,04,11,29,30-02,08"
	sourceLotteryList := []string{
		"DLT:16,18,29,30,31-09,12x3",
		"DLT:01,12,22,30,33-07,10x3",
		"DLT:06,10,21,22,31-07,09x3",
		"DLT:07,10,15,31,33-10,12x3",
		"DLT:08,14,17,27,31-01,07x3",
	}

	target, _ := lottery.GetLottery(targetLottery)

	for _, str := range sourceLotteryList {
		lott, _ := lottery.GetLottery(str)
		result, _ := lott.GetLotteryResult(target)

		result.PrintResult(true, true)
		result.PrintList(true, true)
	}

	dlt.CheckStore()
}
