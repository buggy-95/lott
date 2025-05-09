package lottery

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestParseNumParts(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		errorMsg string
		dan      []int
		tuo      []int
	}{
		{"解析应该成功，1胆1拖", "01~02", "", []int{1}, []int{2}},
		{"解析应该成功，多胆1拖", "01,02~03", "", []int{1, 2}, []int{3}},
		{"解析应该成功，1胆多拖", "01~02,03", "", []int{1}, []int{2, 3}},
		{"解析应该成功，0胆多拖(复试)", "01,02,03", "", nil, []int{1, 2, 3}},
		{"解析应该失败，有错误字符", "01,0-2,03", "号码区解析失败，错误的字符: 【-】。", nil, nil},
		{"解析应该失败，逗号开头", ",01,02,03", "号码至少为1位", nil, nil},
		{"解析应该失败，逗号结尾(有胆码)", "01~02,03,", "号码至少为1位", nil, nil},
		{"解析应该失败，逗号结尾(无胆码)", "01,02,03,", "号码至少为1位", nil, nil},
		{"解析应该失败，波浪号开头", "~01,02,03", "号码至少为1位", nil, nil},
		{"解析应该失败，波浪号结尾", "01,02,03~", "号码至少为1位", nil, nil},
		{"解析应该失败，拖码区重复", "01~02,03~04", "号码区解析错误。拖区重复", nil, nil},
		{"解析应该失败，复试逗号重复", "01,02,,03", "号码至少为1位", nil, nil},
		{"解析应该失败，胆码区逗号重复", "01,,02~03,04", "号码至少为1位", nil, nil},
		{"解析应该失败，拖码区逗号重复", "01,02~03,,04", "号码至少为1位", nil, nil},
		{"解析应该失败，波浪号重复", "01,02~~03,04", "号码至少为1位", nil, nil},
		{"解析应该失败，复试号码过长", "01,002,03", "号码最多为2位数", nil, nil},
		{"解析应该失败，胆码区号码过长", "01,002~03", "号码最多为2位数", nil, nil},
		{"解析应该失败，拖码区号码过长", "01~002,03", "号码最多为2位数", nil, nil},
		{"解析应该失败，复试号码重复", "01,02,03,03,02", "号码区解析失败。拖码区重复: [2 3]。", nil, nil},
		{"解析应该失败，胆码区号码重复", "01,02,02~03,04", "号码区解析失败。胆码区重复: [2]。", nil, nil},
		{"解析应该失败，拖码区号码重复", "01,02~03,04,03", "号码区解析失败。拖码区重复: [3]。", nil, nil},
		{"解析应该失败，胆码区号码重复, 拖码区号码重复", "01,02,01~03,04,03", "号码区解析失败。胆码区重复: [1]。拖码区重复: [3]。", nil, nil},
		{"解析应该失败，拖码区与胆码区冲突", "01,02~02,03", "号码区解析失败。拖码区冲突: [2]。", nil, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dan, tuo, err := parseNumParts(tt.input)

			if len(tt.errorMsg) == 0 {
				if err != nil {
					t.Errorf("%s: 应该解析成功。输入: %s, 错误信息: %s", tt.name, tt.input, err)
				} else if !(reflect.DeepEqual(dan, tt.dan) && reflect.DeepEqual(tuo, tt.tuo)) {
					t.Errorf("%s: 解析结果错误。输入: %s。预期: 胆码 %v, 拖码 %v。实际: 胆码 %v, 拖码 %v", tt.name, tt.input, tt.dan, tt.tuo, dan, tuo)
				}
			} else {
				if err == nil {
					t.Errorf("%s: 应该解析错误。输入: %s", tt.name, tt.input)
				} else if !strings.HasPrefix(err.Error(), tt.errorMsg) {
					t.Errorf("%s: 错误信息错误。预期: %s, 实际: %s", tt.name, tt.errorMsg, err)
				}
			}
		})
	}
}

func TestParseComplexLotteryParts(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		errorMsg string
		parts    ComplexLotteryParts
	}{
		{"解析应该成功，前区胆拖后区复试无倍投无期号", "DLT:01,02,03~04,05-06,07", "", ComplexLotteryParts{LotteryBaseInfo{"DLT", 0, 1}, []int{1, 2, 3}, []int{4, 5}, nil, []int{6, 7}}},
		{"解析应该成功，前区复试后区胆拖无倍投无期号", "DLT:01,02,03,04,05-06~07", "", ComplexLotteryParts{LotteryBaseInfo{"DLT", 0, 1}, nil, []int{1, 2, 3, 4, 5}, []int{6}, []int{7}}},
		{"解析应该成功，前区胆拖后区胆拖无倍投无期号", "DLT:01,02,03~04,05-06~07", "", ComplexLotteryParts{LotteryBaseInfo{"DLT", 0, 1}, []int{1, 2, 3}, []int{4, 5}, []int{6}, []int{7}}},
		{"解析应该成功，复式无倍投无期号", "DLT:01,02,03,04,05-06,07", "", ComplexLotteryParts{LotteryBaseInfo{"DLT", 0, 1}, nil, []int{1, 2, 3, 4, 5}, nil, []int{6, 7}}},
		{"解析应该成功，复式3倍投无期号", "DLT:01,02,03,04,05-06,07x3", "", ComplexLotteryParts{LotteryBaseInfo{"DLT", 0, 3}, nil, []int{1, 2, 3, 4, 5}, nil, []int{6, 7}}},
		{"解析应该成功，复式无倍投有期号", "DLT:01,02,03,04,05-06,07:25053", "", ComplexLotteryParts{LotteryBaseInfo{"DLT", 25053, 1}, nil, []int{1, 2, 3, 4, 5}, nil, []int{6, 7}}},
		{"解析应该成功，复式有期号3倍投", "DLT:01,02,03,04,05-06,07:25053x3", "", ComplexLotteryParts{LotteryBaseInfo{"DLT", 25053, 3}, nil, []int{1, 2, 3, 4, 5}, nil, []int{6, 7}}},
		{"解析应该成功，复式3倍投有期号", "DLT:01,02,03,04,05-06,07x3:25053", "", ComplexLotteryParts{LotteryBaseInfo{"DLT", 25053, 3}, nil, []int{1, 2, 3, 4, 5}, nil, []int{6, 7}}},
		{"解析应该失败，倍投重复", "DLT:01,02,03,04,05-06,07x3:25053x3", "解析失败。倍投已解析过。", ComplexLotteryParts{}},
		{"解析应该失败，连续倍投", "DLT:01,02,03,04,05-06,07x3x3:25053", "倍投数解析失败。当前字符: 【x】", ComplexLotteryParts{}},
		{"解析应该失败，期号重复", "DLT:01,02,03,04,05-06,07:25053x3:25053", "解析失败。期号已解析过。", ComplexLotteryParts{}},
		{"解析应该失败，连续期号", "DLT:01,02,03,04,05-06,07:25053:25053x3", "期数解析失败。当前字符: 【:】。", ComplexLotteryParts{}},
		{"解析应该失败，错误的彩票类型", "DDLT:01,02,03,04,05-06,07x3:25053", "彩票类型解析失败。不支持的彩票类型: DDLT。", ComplexLotteryParts{}},
		{"解析应该失败，前区为空", "DLT:-01,02,03,04,05,06,07x3:25053", "前区号码解析失败。原因: 号码至少为1位。", ComplexLotteryParts{}},
		{"解析应该失败，没有后区", "DLT:01,02,03,04,05,06,07x3:25053", "前区号码解析失败。当前字符: 【x】。", ComplexLotteryParts{}},
		{"解析应该失败，后区为空", "DLT:01,02,03,04,05,06,07-x3:25053", "后区号码解析失败。原因: 号码至少为1位。", ComplexLotteryParts{}},
		{"解析应该失败，倍投错误", "DLT:01,02,03,04,05-06,07xx3:25053", "倍投数解析失败。当前字符: 【x】", ComplexLotteryParts{}},
		{"解析应该失败，倍投错误", "DLT:01,02,03,04,05-06,07x3a:25053", "倍投数解析失败。当前字符: 【a】", ComplexLotteryParts{}},
		{"解析应该失败，期号错误", "DLT:01,02,03,04,05-06,07x3::25053", "期数解析失败。当前字符: 【:】", ComplexLotteryParts{}},
		{"解析应该失败，期号错误", "DLT:01,02,03,04,05-06,07x3:25b053", "期数解析失败。当前字符: 【b】", ComplexLotteryParts{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lottery, err := parseComplexLotteryParts(tt.input)

			if len(tt.errorMsg) == 0 {
				if err != nil {
					t.Errorf("%s: 应该解析成功。输入: %s, 错误信息: %s", tt.name, tt.input, err)
				} else if !reflect.DeepEqual(lottery, tt.parts) {
					t.Errorf("%s: 解析结果错误。输入: %s。预期: %+v, 实际: %+v", tt.name, tt.input, tt.parts, lottery)
				}
			} else {
				if err == nil {
					t.Errorf("%s: 应该解析错误。输入: %s", tt.name, tt.input)
				} else if !strings.HasPrefix(err.Error(), tt.errorMsg) {
					t.Errorf("%s: 错误信息错误。预期: %s, 实际: %s", tt.name, tt.errorMsg, err)
				}
			}
		})
	}
}

func TestGenPermutation(t *testing.T) {
	tests := []struct {
		inputNums []int
		inputN    int
		result    [][]int
	}{
		{[]int{}, 0, [][]int{{}}},
		{[]int{}, 1, [][]int{{}}},
		{[]int{1}, 0, [][]int{{}}},
		{[]int{1}, 1, [][]int{{1}}},
		{[]int{1}, 2, [][]int{{1}}},
		{[]int{1, 2}, 0, [][]int{{}}},
		{[]int{1, 2}, 1, [][]int{{1}, {2}}},
		{[]int{1, 2}, 2, [][]int{{1, 2}}},
		{[]int{1, 2}, 3, [][]int{{1, 2}}},
		{[]int{1, 2, 3}, 1, [][]int{{1}, {2}, {3}}},
		{[]int{1, 2, 3}, 2, [][]int{{1, 2}, {1, 3}, {2, 3}}},
		{[]int{1, 2, 3}, 3, [][]int{{1, 2, 3}}},
		{[]int{1, 2, 3}, 4, [][]int{{1, 2, 3}}},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v %d", tt.inputNums, tt.inputN), func(t *testing.T) {
			result := genPermutation(tt.inputNums, tt.inputN)

			if !reflect.DeepEqual(tt.result, result) {
				t.Errorf("%v %d 失败。预期: %v, 实际: %v", tt.inputNums, tt.inputN, tt.result, result)
			}
		})
	}
}

func TestGetDupNums(t *testing.T) {
	tests := []struct {
		input  []int
		result []int
	}{
		{[]int{}, nil},
		{[]int{1}, nil},
		{[]int{1, 2}, nil},
		{[]int{1, 2, 2}, []int{2}},
		{[]int{1, 2, 2, 1}, []int{1, 2}},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.input), func(t *testing.T) {
			result := GetDupNums(tt.input)

			if !reflect.DeepEqual(tt.result, result) {
				t.Errorf("%v: 期望: %v, 实际: %v", tt.input, tt.result, result)
			}
		})
	}
}

func TestGetCrossNums(t *testing.T) {
	tests := []struct {
		inputSource []int
		inputTarget []int
		result      []int
	}{
		{[]int{}, []int{}, nil},
		{[]int{1}, []int{}, nil},
		{[]int{}, []int{1}, nil},
		{[]int{1}, []int{2}, nil},
		{[]int{1, 2}, []int{2, 3}, []int{2}},
		{[]int{1, 2, 3}, []int{4, 3, 2}, []int{2, 3}},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v%v", tt.inputSource, tt.inputTarget), func(t *testing.T) {
			result := GetCrossNums(tt.inputSource, tt.inputTarget)

			if !reflect.DeepEqual(tt.result, result) {
				t.Errorf("source: %v, target: %v, 期望: %v, 实际: %v", tt.inputSource, tt.inputTarget, tt.result, result)
			}
		})
	}
}

func TestGenSingleLotteryList(t *testing.T) {
	baseInfo := LotteryBaseInfo{"DLT", 0, 1}

	tests := []struct {
		name   string
		input  string
		result []SingleLottery
	}{
		{"单式", "01,02,03,04,05-01,02", []SingleLottery{
			{baseInfo, []int{1, 2, 3, 4, 5}, []int{1, 2}},
		}},
		{"前单后复", "01,02,03,04,05-01,02,03", []SingleLottery{
			{baseInfo, []int{1, 2, 3, 4, 5}, []int{1, 2}},
			{baseInfo, []int{1, 2, 3, 4, 5}, []int{1, 3}},
			{baseInfo, []int{1, 2, 3, 4, 5}, []int{2, 3}},
		}},
		{"前单后拖", "01,02,03,04,05-01~02,03,04", []SingleLottery{
			{baseInfo, []int{1, 2, 3, 4, 5}, []int{1, 2}},
			{baseInfo, []int{1, 2, 3, 4, 5}, []int{1, 3}},
			{baseInfo, []int{1, 2, 3, 4, 5}, []int{1, 4}},
		}},
		{"前复后单", "01,02,03,04,05,06-01,02", []SingleLottery{
			{baseInfo, []int{1, 2, 3, 4, 5}, []int{1, 2}},
			{baseInfo, []int{1, 2, 3, 4, 6}, []int{1, 2}},
			{baseInfo, []int{1, 2, 3, 5, 6}, []int{1, 2}},
			{baseInfo, []int{1, 2, 4, 5, 6}, []int{1, 2}},
			{baseInfo, []int{1, 3, 4, 5, 6}, []int{1, 2}},
			{baseInfo, []int{2, 3, 4, 5, 6}, []int{1, 2}},
		}},
		{"复式", "01,02,03,04,05,06-01,02,03", []SingleLottery{
			{baseInfo, []int{1, 2, 3, 4, 5}, []int{1, 2}},
			{baseInfo, []int{1, 2, 3, 4, 5}, []int{1, 3}},
			{baseInfo, []int{1, 2, 3, 4, 5}, []int{2, 3}},
			{baseInfo, []int{1, 2, 3, 4, 6}, []int{1, 2}},
			{baseInfo, []int{1, 2, 3, 4, 6}, []int{1, 3}},
			{baseInfo, []int{1, 2, 3, 4, 6}, []int{2, 3}},
			{baseInfo, []int{1, 2, 3, 5, 6}, []int{1, 2}},
			{baseInfo, []int{1, 2, 3, 5, 6}, []int{1, 3}},
			{baseInfo, []int{1, 2, 3, 5, 6}, []int{2, 3}},
			{baseInfo, []int{1, 2, 4, 5, 6}, []int{1, 2}},
			{baseInfo, []int{1, 2, 4, 5, 6}, []int{1, 3}},
			{baseInfo, []int{1, 2, 4, 5, 6}, []int{2, 3}},
			{baseInfo, []int{1, 3, 4, 5, 6}, []int{1, 2}},
			{baseInfo, []int{1, 3, 4, 5, 6}, []int{1, 3}},
			{baseInfo, []int{1, 3, 4, 5, 6}, []int{2, 3}},
			{baseInfo, []int{2, 3, 4, 5, 6}, []int{1, 2}},
			{baseInfo, []int{2, 3, 4, 5, 6}, []int{1, 3}},
			{baseInfo, []int{2, 3, 4, 5, 6}, []int{2, 3}},
		}},
		{"前复后拖", "01,02,03,04,05,06-01~02,03,04", []SingleLottery{
			{baseInfo, []int{1, 2, 3, 4, 5}, []int{1, 2}},
			{baseInfo, []int{1, 2, 3, 4, 5}, []int{1, 3}},
			{baseInfo, []int{1, 2, 3, 4, 5}, []int{1, 4}},
			{baseInfo, []int{1, 2, 3, 4, 6}, []int{1, 2}},
			{baseInfo, []int{1, 2, 3, 4, 6}, []int{1, 3}},
			{baseInfo, []int{1, 2, 3, 4, 6}, []int{1, 4}},
			{baseInfo, []int{1, 2, 3, 5, 6}, []int{1, 2}},
			{baseInfo, []int{1, 2, 3, 5, 6}, []int{1, 3}},
			{baseInfo, []int{1, 2, 3, 5, 6}, []int{1, 4}},
			{baseInfo, []int{1, 2, 4, 5, 6}, []int{1, 2}},
			{baseInfo, []int{1, 2, 4, 5, 6}, []int{1, 3}},
			{baseInfo, []int{1, 2, 4, 5, 6}, []int{1, 4}},
			{baseInfo, []int{1, 3, 4, 5, 6}, []int{1, 2}},
			{baseInfo, []int{1, 3, 4, 5, 6}, []int{1, 3}},
			{baseInfo, []int{1, 3, 4, 5, 6}, []int{1, 4}},
			{baseInfo, []int{2, 3, 4, 5, 6}, []int{1, 2}},
			{baseInfo, []int{2, 3, 4, 5, 6}, []int{1, 3}},
			{baseInfo, []int{2, 3, 4, 5, 6}, []int{1, 4}},
		}},
		{"前拖后单", "01,02,03,04~05,06-01,02", []SingleLottery{
			{baseInfo, []int{1, 2, 3, 4, 5}, []int{1, 2}},
			{baseInfo, []int{1, 2, 3, 4, 6}, []int{1, 2}},
		}},
		{"前拖后复", "01,02,03,04~05,06-01,02,03", []SingleLottery{
			{baseInfo, []int{1, 2, 3, 4, 5}, []int{1, 2}},
			{baseInfo, []int{1, 2, 3, 4, 5}, []int{1, 3}},
			{baseInfo, []int{1, 2, 3, 4, 5}, []int{2, 3}},
			{baseInfo, []int{1, 2, 3, 4, 6}, []int{1, 2}},
			{baseInfo, []int{1, 2, 3, 4, 6}, []int{1, 3}},
			{baseInfo, []int{1, 2, 3, 4, 6}, []int{2, 3}},
		}},
		{"前拖后拖", "01,02,03,04~05,06-01~02,03,04", []SingleLottery{
			{baseInfo, []int{1, 2, 3, 4, 5}, []int{1, 2}},
			{baseInfo, []int{1, 2, 3, 4, 5}, []int{1, 3}},
			{baseInfo, []int{1, 2, 3, 4, 5}, []int{1, 4}},
			{baseInfo, []int{1, 2, 3, 4, 6}, []int{1, 2}},
			{baseInfo, []int{1, 2, 3, 4, 6}, []int{1, 3}},
			{baseInfo, []int{1, 2, 3, 4, 6}, []int{1, 4}},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			complexLottery, err := parseComplexLotteryParts("DLT:" + tt.input)

			if err != nil {
				t.Errorf("%s: 解析失败。错误信息: %s, 输入: %s", tt.name, err, tt.input)

				return
			}

			result := genSingleLotteryList(complexLottery)

			if !reflect.DeepEqual(tt.result, result) {
				t.Errorf("%s。 期望: %v, 实际: %v, 输入: %+v", tt.name, tt.result, result, tt.input)
			}
		})
	}
}

func TestGetMatchResult(t *testing.T) {
	tests := []struct {
		name    string
		source  []int
		target  []int
		result  []BingoNum
		matched int
	}{
		{"source和target都为空", []int{}, []int{}, nil, 0},
		{"source为空", []int{}, []int{1, 2}, nil, 0},
		{"target为空", []int{1, 2}, []int{}, []BingoNum{{1, false}, {2, false}}, 0},
		{"source和target一样", []int{1, 2}, []int{2, 1}, []BingoNum{{1, true}, {2, true}}, 2},
		{"source和target有交集", []int{1, 2}, []int{2, 3}, []BingoNum{{1, false}, {2, true}}, 1},
		{"source和target没有交集", []int{1, 2}, []int{3, 4}, []BingoNum{{1, false}, {2, false}}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, matched := getMatchNums(tt.source, tt.target)

			if !reflect.DeepEqual(tt.result, result) && matched != tt.matched {
				t.Errorf("%s: 期望: %v %d, 实际: %v %d", tt.name, tt.result, tt.matched, result, matched)
			}
		})
	}
}

func TestGetSingleLotteryResult(t *testing.T) {
	baseInfo := LotteryBaseInfo{"DLT", 0, 1}

	tests := []struct {
		name   string
		source string
		target string
		result SingleLotteryResult
	}{
		{"一等奖", "01,02,03,04,05-01,02", "01,02,03,04,05-01,02", SingleLotteryResult{baseInfo, 5, 2, 1, 10000000, []ResultNum{
			{BingoNum{1, true}, "FrontTuo"},
			{BingoNum{2, true}, "FrontTuo"},
			{BingoNum{3, true}, "FrontTuo"},
			{BingoNum{4, true}, "FrontTuo"},
			{BingoNum{5, true}, "FrontTuo"},
			{BingoNum{1, true}, "BackTuo"},
			{BingoNum{2, true}, "BackTuo"},
		}}},
		{"二等奖", "01,02,03,04,05-01,03", "01,02,03,04,05-01,02", SingleLotteryResult{baseInfo, 5, 1, 2, 200000, []ResultNum{
			{BingoNum{1, true}, "FrontTuo"},
			{BingoNum{2, true}, "FrontTuo"},
			{BingoNum{3, true}, "FrontTuo"},
			{BingoNum{4, true}, "FrontTuo"},
			{BingoNum{5, true}, "FrontTuo"},
			{BingoNum{1, true}, "BackTuo"},
			{BingoNum{3, false}, "BackTuo"},
		}}},
		{"三等奖", "01,02,03,04,05-03,04", "01,02,03,04,05-01,02", SingleLotteryResult{baseInfo, 5, 0, 3, 10000, []ResultNum{
			{BingoNum{1, true}, "FrontTuo"},
			{BingoNum{2, true}, "FrontTuo"},
			{BingoNum{3, true}, "FrontTuo"},
			{BingoNum{4, true}, "FrontTuo"},
			{BingoNum{5, true}, "FrontTuo"},
			{BingoNum{3, false}, "BackTuo"},
			{BingoNum{4, false}, "BackTuo"},
		}}},
		{"四等奖", "01,02,03,04,06-01,02", "01,02,03,04,05-01,02", SingleLotteryResult{baseInfo, 4, 2, 4, 3000, []ResultNum{
			{BingoNum{1, true}, "FrontTuo"},
			{BingoNum{2, true}, "FrontTuo"},
			{BingoNum{3, true}, "FrontTuo"},
			{BingoNum{4, true}, "FrontTuo"},
			{BingoNum{6, false}, "FrontTuo"},
			{BingoNum{1, true}, "BackTuo"},
			{BingoNum{2, true}, "BackTuo"},
		}}},
		{"五等奖", "01,02,03,04,06-01,03", "01,02,03,04,05-01,02", SingleLotteryResult{baseInfo, 4, 1, 5, 300, []ResultNum{
			{BingoNum{1, true}, "FrontTuo"},
			{BingoNum{2, true}, "FrontTuo"},
			{BingoNum{3, true}, "FrontTuo"},
			{BingoNum{4, true}, "FrontTuo"},
			{BingoNum{6, false}, "FrontTuo"},
			{BingoNum{1, true}, "BackTuo"},
			{BingoNum{3, false}, "BackTuo"},
		}}},
		{"六等奖", "01,02,03,06,07-01,02", "01,02,03,04,05-01,02", SingleLotteryResult{baseInfo, 3, 2, 6, 200, []ResultNum{
			{BingoNum{1, true}, "FrontTuo"},
			{BingoNum{2, true}, "FrontTuo"},
			{BingoNum{3, true}, "FrontTuo"},
			{BingoNum{6, false}, "FrontTuo"},
			{BingoNum{7, false}, "FrontTuo"},
			{BingoNum{1, true}, "BackTuo"},
			{BingoNum{2, true}, "BackTuo"},
		}}},
		{"七等奖", "01,02,03,04,06-03,04", "01,02,03,04,05-01,02", SingleLotteryResult{baseInfo, 4, 0, 7, 100, []ResultNum{
			{BingoNum{1, true}, "FrontTuo"},
			{BingoNum{2, true}, "FrontTuo"},
			{BingoNum{3, true}, "FrontTuo"},
			{BingoNum{4, true}, "FrontTuo"},
			{BingoNum{6, false}, "FrontTuo"},
			{BingoNum{3, false}, "BackTuo"},
			{BingoNum{4, false}, "BackTuo"},
		}}},
		{"八等奖A", "01,02,03,06,07-01,03", "01,02,03,04,05-01,02", SingleLotteryResult{baseInfo, 3, 1, 8, 15, []ResultNum{
			{BingoNum{1, true}, "FrontTuo"},
			{BingoNum{2, true}, "FrontTuo"},
			{BingoNum{3, true}, "FrontTuo"},
			{BingoNum{6, false}, "FrontTuo"},
			{BingoNum{7, false}, "FrontTuo"},
			{BingoNum{1, true}, "BackTuo"},
			{BingoNum{3, false}, "BackTuo"},
		}}},
		{"八等奖B", "01,02,06,07,08-01,02", "01,02,03,04,05-01,02", SingleLotteryResult{baseInfo, 2, 2, 8, 15, []ResultNum{
			{BingoNum{1, true}, "FrontTuo"},
			{BingoNum{2, true}, "FrontTuo"},
			{BingoNum{6, false}, "FrontTuo"},
			{BingoNum{7, false}, "FrontTuo"},
			{BingoNum{8, false}, "FrontTuo"},
			{BingoNum{1, true}, "BackTuo"},
			{BingoNum{2, true}, "BackTuo"},
		}}},
		{"九等奖A", "01,02,03,06,07-03,04", "01,02,03,04,05-01,02", SingleLotteryResult{baseInfo, 3, 0, 9, 5, []ResultNum{
			{BingoNum{1, true}, "FrontTuo"},
			{BingoNum{2, true}, "FrontTuo"},
			{BingoNum{3, true}, "FrontTuo"},
			{BingoNum{6, false}, "FrontTuo"},
			{BingoNum{7, false}, "FrontTuo"},
			{BingoNum{3, false}, "BackTuo"},
			{BingoNum{4, false}, "BackTuo"},
		}}},
		{"九等奖B", "01,06,07,08,09-01,02", "01,02,03,04,05-01,02", SingleLotteryResult{baseInfo, 1, 2, 9, 5, []ResultNum{
			{BingoNum{1, true}, "FrontTuo"},
			{BingoNum{6, false}, "FrontTuo"},
			{BingoNum{7, false}, "FrontTuo"},
			{BingoNum{8, false}, "FrontTuo"},
			{BingoNum{9, false}, "FrontTuo"},
			{BingoNum{1, true}, "BackTuo"},
			{BingoNum{2, true}, "BackTuo"},
		}}},
		{"九等奖C", "01,02,06,07,08-01,03", "01,02,03,04,05-01,02", SingleLotteryResult{baseInfo, 2, 1, 9, 5, []ResultNum{
			{BingoNum{1, true}, "FrontTuo"},
			{BingoNum{2, true}, "FrontTuo"},
			{BingoNum{6, false}, "FrontTuo"},
			{BingoNum{7, false}, "FrontTuo"},
			{BingoNum{8, false}, "FrontTuo"},
			{BingoNum{1, true}, "BackTuo"},
			{BingoNum{3, false}, "BackTuo"},
		}}},
		{"九等奖D", "06,07,08,09,10-01,02", "01,02,03,04,05-01,02", SingleLotteryResult{baseInfo, 0, 2, 9, 5, []ResultNum{
			{BingoNum{6, false}, "FrontTuo"},
			{BingoNum{7, false}, "FrontTuo"},
			{BingoNum{8, false}, "FrontTuo"},
			{BingoNum{9, false}, "FrontTuo"},
			{BingoNum{10, false}, "FrontTuo"},
			{BingoNum{1, true}, "BackTuo"},
			{BingoNum{2, true}, "BackTuo"},
		}}},
		{"无奖A", "06,07,08,09,10-03,04", "01,02,03,04,05-01,02", SingleLotteryResult{baseInfo, 0, 0, 0, 0, []ResultNum{
			{BingoNum{6, false}, "FrontTuo"},
			{BingoNum{7, false}, "FrontTuo"},
			{BingoNum{8, false}, "FrontTuo"},
			{BingoNum{9, false}, "FrontTuo"},
			{BingoNum{10, false}, "FrontTuo"},
			{BingoNum{3, false}, "BackTuo"},
			{BingoNum{4, false}, "BackTuo"},
		}}},
		{"无奖B", "01,06,07,08,09-03,04", "01,02,03,04,05-01,02", SingleLotteryResult{baseInfo, 1, 0, 0, 0, []ResultNum{
			{BingoNum{1, true}, "FrontTuo"},
			{BingoNum{6, false}, "FrontTuo"},
			{BingoNum{7, false}, "FrontTuo"},
			{BingoNum{8, false}, "FrontTuo"},
			{BingoNum{9, false}, "FrontTuo"},
			{BingoNum{3, false}, "BackTuo"},
			{BingoNum{4, false}, "BackTuo"},
		}}},
		{"无奖C", "06,07,08,09,10-01,03", "01,02,03,04,05-01,02", SingleLotteryResult{baseInfo, 0, 1, 0, 0, []ResultNum{
			{BingoNum{6, false}, "FrontTuo"},
			{BingoNum{7, false}, "FrontTuo"},
			{BingoNum{8, false}, "FrontTuo"},
			{BingoNum{9, false}, "FrontTuo"},
			{BingoNum{10, false}, "FrontTuo"},
			{BingoNum{1, true}, "BackTuo"},
			{BingoNum{3, false}, "BackTuo"},
		}}},
		{"无奖D", "01,06,07,08,09-01,03", "01,02,03,04,05-01,02", SingleLotteryResult{baseInfo, 1, 1, 0, 0, []ResultNum{
			{BingoNum{1, true}, "FrontTuo"},
			{BingoNum{6, false}, "FrontTuo"},
			{BingoNum{7, false}, "FrontTuo"},
			{BingoNum{8, false}, "FrontTuo"},
			{BingoNum{9, false}, "FrontTuo"},
			{BingoNum{1, true}, "BackTuo"},
			{BingoNum{3, false}, "BackTuo"},
		}}},
		{"无奖E", "01,02,06,07,08-03,04", "01,02,03,04,05-01,02", SingleLotteryResult{baseInfo, 2, 0, 0, 0, []ResultNum{
			{BingoNum{1, true}, "FrontTuo"},
			{BingoNum{2, true}, "FrontTuo"},
			{BingoNum{6, false}, "FrontTuo"},
			{BingoNum{7, false}, "FrontTuo"},
			{BingoNum{8, false}, "FrontTuo"},
			{BingoNum{3, false}, "BackTuo"},
			{BingoNum{4, false}, "BackTuo"},
		}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			complexSource, _ := parseComplexLotteryParts("DLT:" + tt.source)
			sourceList := genSingleLotteryList(complexSource)
			source := sourceList[0]
			complexTarget, _ := parseComplexLotteryParts("DLT:" + tt.target)
			targetList := genSingleLotteryList(complexTarget)
			target := targetList[0]
			result := getSingleLotteryResult(source, target)

			if !reflect.DeepEqual(tt.result, result) {
				t.Errorf("期望: %v, 实际: %v, source: %+v, target: %+v", tt.result, result, tt.source, tt.target)
			}
		})
	}
}

func TestFormatSingleLottery(t *testing.T) {
	tests := []struct {
		name        string
		onlyNumbers bool
		input       string
		result      string
	}{
		{"无倍投无期号全展示", false, "DLT:01,02,03,04,05-01,02", "01,02,03,04,05-01,02"},
		{"3倍投无期号全展示", false, "DLT:01,02,03,04,05-01,02x3", "01,02,03,04,05-01,02x3"},
		{"无倍投有期号全展示", false, "DLT:01,02,03,04,05-01,02:25053", "01,02,03,04,05-01,02:25053"},
		{"3倍投有期号全展示", false, "DLT:01,02,03,04,05-01,02:25053x3", "01,02,03,04,05-01,02x3:25053"},
		{"无倍投无期号仅数字", true, "DLT:01,02,03,04,05-01,02", "01,02,03,04,05-01,02"},
		{"3倍投无期号仅数字", true, "DLT:01,02,03,04,05-01,02x3", "01,02,03,04,05-01,02"},
		{"无倍投有期号仅数字", true, "DLT:01,02,03,04,05-01,02:25053", "01,02,03,04,05-01,02"},
		{"3倍投有期号仅数字", true, "DLT:01,02,03,04,05-01,02:25053x3", "01,02,03,04,05-01,02"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			complexLottery, err := GetComplexLottery(tt.input)

			if err != nil {
				t.Errorf("%s: 解析失败。错误信息: %s, 输入: %s", tt.name, err, tt.input)
				return
			}

			result := formatSingleLottery(complexLottery.List[0], tt.onlyNumbers)

			if result != tt.result {
				t.Errorf("%s。期望: %s, 实际: %s, 输入: %s %+v", tt.name, tt.result, result, tt.input, complexLottery)
			}
		})
	}
}
