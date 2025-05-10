package dlt

import (
	"errors"
	"fmt"

	"github.com/buggy-95/lott/internal/lottery"
)

// check
//
// @Description 检查复式彩票的前区和后区是否符合要求
//
// @Param lott ComplexLotteryParts 复杂彩票结构体
//
// @Return error 错误信息
func check(lott lottery.LotteryParts) error {
	if len(lott.FrontDan) >= 5 {
		return fmt.Errorf("前区胆码数量应该小于5，当前数量: %d", len(lott.FrontDan))
	} else if len(lott.BackDan) >= 2 {
		return fmt.Errorf("后区胆码数量应该小于2，当前数量: %d", len(lott.BackDan))
	} else if arr := lottery.GetDupNums(lott.FrontDan); len(arr) > 0 {
		return fmt.Errorf("前区胆码重复: %v", arr)
	} else if arr := lottery.GetDupNums(lott.FrontTuo); len(arr) > 0 {
		return fmt.Errorf("前区拖码重复: %v", arr)
	} else if arr := lottery.GetDupNums(lott.BackTuo); len(arr) > 0 {
		return fmt.Errorf("后区拖码重复: %v", arr)
	} else if arr := lottery.GetCrossNums(lott.FrontDan, lott.FrontTuo); len(arr) > 0 {
		return fmt.Errorf("前区拖码与胆码重复: %v", arr)
	} else if arr := lottery.GetCrossNums(lott.BackDan, lott.BackTuo); len(arr) > 0 {
		return fmt.Errorf("后区拖码与胆码重复: %v", arr)
	}

	front := append(lott.FrontDan, lott.FrontTuo...)
	back := append(lott.BackDan, lott.BackTuo...)

	if len(front) < 5 {
		return errors.New("前区最少需要5个数字")
	}

	if len(back) < 2 {
		return errors.New("后区最少需要2个数字")
	}

	for _, n := range front {
		if !(1 <= n && n <= 35) {
			return errors.New("前区数字范围为1~35")
		}
	}

	for _, n := range back {
		if !(1 <= n && n <= 12) {
			return errors.New("后区数字范围为1~12")
		}
	}

	return nil
}
