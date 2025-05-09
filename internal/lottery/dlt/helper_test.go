package dlt

import (
	"testing"

	"github.com/buggy-95/lott/internal/lottery"
)

func TestCheck(t *testing.T) {
	tests := []struct {
		name  string
		msg   string
		input lottery.ComplexLotteryParts
	}{
		{"应该成功，单式", "", lottery.ComplexLotteryParts{FrontDan: []int{}, FrontTuo: []int{1, 2, 3, 4, 5}, BackDan: []int{}, BackTuo: []int{1, 2}}},
		{"应该成功，前复后单", "", lottery.ComplexLotteryParts{FrontDan: []int{}, FrontTuo: []int{1, 2, 3, 4, 5, 6}, BackDan: []int{}, BackTuo: []int{1, 2}}},
		{"应该成功，前单后复", "", lottery.ComplexLotteryParts{FrontDan: []int{}, FrontTuo: []int{1, 2, 3, 4, 5}, BackDan: []int{}, BackTuo: []int{1, 2, 3}}},
		{"应该成功，复式", "", lottery.ComplexLotteryParts{FrontDan: []int{}, FrontTuo: []int{1, 2, 3, 4, 5, 6}, BackDan: []int{}, BackTuo: []int{1, 2, 3}}},
		{"应该成功，前拖", "", lottery.ComplexLotteryParts{FrontDan: []int{1, 2, 3}, FrontTuo: []int{4, 5, 6}, BackDan: []int{}, BackTuo: []int{1, 2}}},
		{"应该成功，后拖", "", lottery.ComplexLotteryParts{FrontDan: []int{}, FrontTuo: []int{1, 2, 3, 4, 5}, BackDan: []int{1}, BackTuo: []int{2, 3}}},
		{"应该成功，前拖后拖", "", lottery.ComplexLotteryParts{FrontDan: []int{1, 2, 3}, FrontTuo: []int{4, 5, 6}, BackDan: []int{1}, BackTuo: []int{2, 3}}},
		{"前区胆码过多", "前区胆码数量应该小于5，当前数量: 6", lottery.ComplexLotteryParts{FrontDan: []int{1, 2, 3, 4, 5, 6}, FrontTuo: []int{7}, BackDan: []int{1}, BackTuo: []int{2}}},
		{"后区胆码过多", "后区胆码数量应该小于2，当前数量: 2", lottery.ComplexLotteryParts{FrontDan: []int{1, 2, 3, 4}, FrontTuo: []int{5, 6, 7}, BackDan: []int{1, 2}, BackTuo: []int{3}}},
		{"前区胆码重复", "前区胆码重复: [1 2]", lottery.ComplexLotteryParts{FrontDan: []int{1, 2, 2, 1}, FrontTuo: []int{5, 6, 7}, BackDan: []int{1}, BackTuo: []int{2, 3}}},
		{"前区拖码重复", "前区拖码重复: [5]", lottery.ComplexLotteryParts{FrontDan: []int{1, 2, 3, 4}, FrontTuo: []int{5, 6, 5}, BackDan: []int{1}, BackTuo: []int{2, 3}}},
		{"前区胆拖交叉", "前区拖码与胆码重复: [3 4]", lottery.ComplexLotteryParts{FrontDan: []int{1, 2, 3, 4}, FrontTuo: []int{5, 4, 3}, BackDan: []int{1}, BackTuo: []int{2, 3}}},
		{"后区胆拖交叉", "后区拖码与胆码重复: [1]", lottery.ComplexLotteryParts{FrontDan: []int{1, 2, 3, 4}, FrontTuo: []int{5, 6, 7}, BackDan: []int{1}, BackTuo: []int{2, 1}}},
		{"前区太少，无拖码", "前区最少需要5个数字", lottery.ComplexLotteryParts{FrontDan: []int{1, 2, 3, 4}, FrontTuo: []int{}, BackDan: []int{1}, BackTuo: []int{2, 3}}},
		{"前区太少，无胆码", "前区最少需要5个数字", lottery.ComplexLotteryParts{FrontDan: []int{}, FrontTuo: []int{1, 2, 3, 4}, BackDan: []int{1}, BackTuo: []int{2, 3}}},
		{"前区太少，胆拖", "前区最少需要5个数字", lottery.ComplexLotteryParts{FrontDan: []int{1, 2}, FrontTuo: []int{3, 4}, BackDan: []int{1}, BackTuo: []int{2, 3}}},
		{"后区太少，无拖码", "后区最少需要2个数字", lottery.ComplexLotteryParts{FrontDan: []int{1, 2, 3}, FrontTuo: []int{4, 5, 6}, BackDan: []int{1}, BackTuo: []int{}}},
		{"后区太少，无胆码", "后区最少需要2个数字", lottery.ComplexLotteryParts{FrontDan: []int{1, 2, 3}, FrontTuo: []int{4, 5, 6}, BackDan: []int{}, BackTuo: []int{1}}},
		{"前区胆码过大", "前区数字范围为1~35", lottery.ComplexLotteryParts{FrontDan: []int{1, 2, 3, 36}, FrontTuo: []int{4, 5, 6}, BackDan: []int{1}, BackTuo: []int{2, 3}}},
		{"前区胆码过小", "前区数字范围为1~35", lottery.ComplexLotteryParts{FrontDan: []int{1, 2, 3, 0}, FrontTuo: []int{4, 5, 6}, BackDan: []int{1}, BackTuo: []int{2, 3}}},
		{"前区拖码过大", "前区数字范围为1~35", lottery.ComplexLotteryParts{FrontDan: []int{1, 2, 3}, FrontTuo: []int{4, 5, 6, 36}, BackDan: []int{1}, BackTuo: []int{2, 3}}},
		{"前区拖码过小", "前区数字范围为1~35", lottery.ComplexLotteryParts{FrontDan: []int{1, 2, 3}, FrontTuo: []int{4, 5, 6, 0}, BackDan: []int{1}, BackTuo: []int{2, 3}}},
		{"后区胆码过大", "后区数字范围为1~12", lottery.ComplexLotteryParts{FrontDan: []int{1, 2, 3}, FrontTuo: []int{4, 5, 6}, BackDan: []int{13}, BackTuo: []int{2, 3}}},
		{"后区胆码过小", "后区数字范围为1~12", lottery.ComplexLotteryParts{FrontDan: []int{1, 2, 3}, FrontTuo: []int{4, 5, 6}, BackDan: []int{0}, BackTuo: []int{2, 3}}},
		{"后区拖码过大", "后区数字范围为1~12", lottery.ComplexLotteryParts{FrontDan: []int{1, 2, 3}, FrontTuo: []int{4, 5, 6}, BackDan: []int{1}, BackTuo: []int{2, 13}}},
		{"后区拖码过小", "后区数字范围为1~12", lottery.ComplexLotteryParts{FrontDan: []int{1, 2, 3}, FrontTuo: []int{4, 5, 6}, BackDan: []int{1}, BackTuo: []int{2, 0}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := check(tt.input)

			if err != nil {
				if len(tt.msg) == 0 {
					t.Errorf("应该成功，错误信息: %s", err)
				} else if err.Error() != tt.msg {
					t.Errorf("错误信息错误，期望: %s, 实际: %s。输入: %+v", tt.msg, err, tt.input)
				}
			} else if len(tt.msg) > 0 {
				t.Errorf("应该失败，输入: %+v", tt.input)
			}
		})
	}
}

func TestIsSingleComplexLottery(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		result bool
	}{
		{"单式", "01,02,03,04,05-01,02", true},
		{"前单后复", "01,02,03,04,05-01,02,03", false},
		{"前单后拖", "01,02,03,04,05-01~02,03", false},
		{"前复后单", "01,02,03,04,05,06-01,02", false},
		{"复式", "01,02,03,04,05,06-01,02,03", false},
		{"前复后拖", "01,02,03,04,05,06-01~02,03", false},
		{"前拖后单", "01,02,03~04,05,06-01,02", false},
		{"前拖后复", "01,02,03~04,05,06-01,02,03", false},
		{"前拖后拖", "01,02,03~04,05,06-01~02,03", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			complexLottery, err := lottery.GetComplexLottery("DLT:" + tt.input)

			if err != nil {
				t.Errorf("解析失败，错误信息: %s", err)
				return
			} else if result := isSingleComplexLottery(complexLottery); result != tt.result {
				t.Errorf("期望: %v, 实际: %v, 输入: %s", tt.result, result, tt.input)
			}
		})
	}
}
