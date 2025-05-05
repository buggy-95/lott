package lottery

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseLottery(t *testing.T) {
	t.Skip()
	defaultLottery := Lottery{
		Type:  "unknown",
		Index: 0,
		Red:   []int{},
		Scale: 1,
		Blue:  []int{},
	}

	type TestCase struct {
		isPass bool
		name   string
		input  string
		result Lottery
	}

	// 仅有号码
	testsOnlyNumbers := []TestCase{
		{false, "解析应该失败-红球超2位", "001,02,03,04,05-06,07", defaultLottery},
		{false, "解析应该失败-蓝球超2位", "01,02,03,04,05-006,07", defaultLottery},
		{false, "解析应该失败-存在错误字符", "01,02,03,04,a05-06,07", defaultLottery},
		{false, "解析应该失败-分隔符连用(,,)", "01,,02,03,04,05-06,07", defaultLottery},
		{false, "解析应该失败-分隔符连用(,-)", "01,02,03,04,-06,07", defaultLottery},
		{false, "解析应该失败-分隔符连用(-,)", "01,02,03,04,05-,07", defaultLottery},
		{false, "解析应该失败-无蓝球", "01,02,03,04,05,06,07", defaultLottery},
		{false, "解析应该失败-逗号开头", ",02,03,04,05-06,07", defaultLottery},
		{false, "解析应该失败-分割号开头", "-02,03,04,05-06,07", defaultLottery},
		{false, "解析应该失败-逗号结尾，无蓝球", "01,02,03,04,05,06,", defaultLottery},
		{false, "解析应该失败-逗号结尾，有蓝球", "01,02,03,04,05-06,", defaultLottery},
		{false, "解析应该失败-分割号结尾，无蓝球", "01,02,03,04,05-", defaultLottery},
		{false, "解析应该失败-分割号结尾，有蓝球", "01,02,03,04,05-06-", defaultLottery},
		{false, "解析应该失败-蓝球区重复", "01,02,03,04,05-06-07", defaultLottery},
		{true, "解析应该成功-正确的SSQ", "01,02,03,04,05,06-07", Lottery{"unknown", 0, 1, []int{1, 2, 3, 4, 5, 6}, []int{7}}},
		{true, "解析应该成功-正确的DLT", "01,02,03,04,05-06,07", Lottery{"unknown", 0, 1, []int{1, 2, 3, 4, 5}, []int{6, 7}}},
	}

	// 包含倍投
	testsWithScale := []TestCase{
		{false, "解析应该失败-存在错误字符", "01,02,03,04,a05-06,07x3", defaultLottery},
		{false, "解析应该失败-分隔符连用(,,)", "01,,02,03,04,05-06,07x3", defaultLottery},
		{false, "解析应该失败-分隔符连用(,-)", "01,02,03,04,-06,07x3", defaultLottery},
		{false, "解析应该失败-分隔符连用(-,)", "01,02,03,04,05-,07x3", defaultLottery},
		{false, "解析应该失败-分隔符连用(,x)", "01,02,03,04,05,06,x3", defaultLottery},
		{false, "解析应该失败-分隔符连用(-x)", "01,02,03,04,05,06-x3", defaultLottery},
		{false, "解析应该失败-分隔符连用(xx)", "01,02,03,04,05,06xx3", defaultLottery},
		{false, "解析应该失败-无蓝球", "01,02,03,04,05,06,07x3", defaultLottery},
		{false, "解析应该失败-逗号开头", ",02,03,04,05-06,07x3", defaultLottery},
		{false, "解析应该失败-分割号开头", "-02,03,04,05-06,07x3", defaultLottery},
		{false, "解析应该失败-x开头", "x3,01,02,03,04,05-06", defaultLottery},
		{false, "解析应该失败-其他字符开头", "a02,03,04,05-06,07x3", defaultLottery},
		{false, "解析应该失败-逗号结尾，无蓝球", "01,02,03,04,05,06x3,", defaultLottery},
		{false, "解析应该失败-逗号结尾，有蓝球", "01,02,03,04,05-06x3,", defaultLottery},
		{false, "解析应该失败-分割号结尾，无蓝球", "01,02,03,04,05x3-", defaultLottery},
		{false, "解析应该失败-分割号结尾，有蓝球", "01,02,03,04,05-06x3-", defaultLottery},
		{false, "解析应该失败-x结尾，无蓝球", "01,02,03,04,05x", defaultLottery},
		{false, "解析应该失败-x结尾，有蓝球", "01,02,03,04,05-06x", defaultLottery},
		{false, "解析应该失败-倍投在中间", "01,02,03x3,04,05-06", defaultLottery},
		{false, "解析应该失败-蓝球区重复", "01,02,03,04,05-06-07x3", defaultLottery},
		{false, "解析应该失败-倍投区重复", "01,02,03,04,05-06,07x3x3", defaultLottery},
		{true, "解析应该成功-正确的SSQ", "01,02,03,04,05,06-07x3", Lottery{"unknown", 0, 3, []int{1, 2, 3, 4, 5, 6}, []int{7}}},
		{true, "解析应该成功-正确的DLT", "01,02,03,04,05-06,07x3", Lottery{"unknown", 0, 3, []int{1, 2, 3, 4, 5}, []int{6, 7}}},
	}

	// 包含期号
	testsWithIndex := []TestCase{
		{false, "解析应该失败-存在错误字符", "01,02,03,04,a05-06,07:12345", defaultLottery},
		{false, "解析应该失败-分隔符连用(,,)", "01,,02,03,04,05-06,07:12345", defaultLottery},
		{false, "解析应该失败-分隔符连用(,-)", "01,02,03,04,-06,07:12345", defaultLottery},
		{false, "解析应该失败-分隔符连用(-,)", "01,02,03,04,05-,07:12345", defaultLottery},
		{false, "解析应该失败-分隔符连用(,:)", "01,02,03,04,05,06,:12345", defaultLottery},
		{false, "解析应该失败-分隔符连用(-:)", "01,02,03,04,05,06-:12345", defaultLottery},
		{false, "解析应该失败-分隔符连用(::)", "01,02,03,04,05,06::12345", defaultLottery},
		{false, "解析应该失败-无蓝球", "01,02,03,04,05,06,07:12345", defaultLottery},
		{false, "解析应该失败-逗号开头", ",02,03,04,05-06,07:12345", defaultLottery},
		{false, "解析应该失败-分割号开头", "-02,03,04,05-06,07:12345", defaultLottery},
		{false, "解析应该失败-x开头", "x3,01,02,03,04,05-06:12345", defaultLottery},
		{false, "解析应该失败-:开头", ":12345,01,02,03,04,05-06:12345", defaultLottery},
		{false, "解析应该失败-其他字符开头", "a02,03,04,05-06,07:12345", defaultLottery},
		{false, "解析应该失败-逗号结尾，无蓝球", "01,02,03,04,05,06:12345,", defaultLottery},
		{false, "解析应该失败-逗号结尾，有蓝球", "01,02,03,04,05-06:12345,", defaultLottery},
		{false, "解析应该失败-分割号结尾，无蓝球", "01,02,03,04,05:12345-", defaultLottery},
		{false, "解析应该失败-分割号结尾，有蓝球", "01,02,03,04,05-06:12345-", defaultLottery},
		{false, "解析应该失败-x结尾，无蓝球", "01,02,03,04,05:12345:", defaultLottery},
		{false, "解析应该失败-x结尾，有蓝球", "01,02,03,04,05-06:12345x", defaultLottery},
		{false, "解析应该失败-:结尾，无蓝球", "01,02,03,04,05:", defaultLottery},
		{false, "解析应该失败-:结尾，有蓝球", "01,02,03,04,05-06:", defaultLottery},
		{false, "解析应该失败-期号在中间", "01,02,03:12345,04,05-06", defaultLottery},
		{false, "解析应该失败-蓝球区重复", "01,02,03,04,05-06-07:12345", defaultLottery},
		{false, "解析应该失败-倍投区重复", "01,02,03,04,05-06,07:12345:12345", defaultLottery},
		{true, "解析应该成功-正确的SSQ", "01,02,03,04,05,06-07:12345", Lottery{"unknown", 12345, 1, []int{1, 2, 3, 4, 5, 6}, []int{7}}},
		{true, "解析应该成功-正确的DLT", "01,02,03,04,05-06,07:12345", Lottery{"unknown", 12345, 1, []int{1, 2, 3, 4, 5}, []int{6, 7}}},
	}

	// 包含倍投和期号
	testsWithScaleAndIndex := []TestCase{
		{false, "解析应该失败-存在错误字符", "01,02,03,04,a05-06,07x3:12345", defaultLottery},
		{false, "解析应该失败-分隔符连用(,,)", "01,,02,03,04,05-06,07x3:12345", defaultLottery},
		{false, "解析应该失败-分隔符连用(,-)", "01,02,03,04,-06,07x3:12345", defaultLottery},
		{false, "解析应该失败-分隔符连用(-,)", "01,02,03,04,05-,07x3:12345", defaultLottery},
		{false, "解析应该失败-分隔符连用(,:)", "01,02,03,04,05,06x3,:12345", defaultLottery},
		{false, "解析应该失败-分隔符连用(-:)", "01,02,03,04,05,06x3-:12345", defaultLottery},
		{false, "解析应该失败-分隔符连用(::)", "01,02,03,04,05,06x3::12345", defaultLottery},
		{false, "解析应该失败-无蓝球", "01,02,03,04,05,06,07x3:12345", defaultLottery},
		{false, "解析应该失败-逗号开头", ",02,03,04,05-06,07x3:12345", defaultLottery},
		{false, "解析应该失败-分割号开头", "-02,03,04,05-06,07x3:12345", defaultLottery},
		{false, "解析应该失败-x开头", "x3,01,02,03,04,05-06x3:12345", defaultLottery},
		{false, "解析应该失败-:开头", ":12345,01,02,03,04,05-06x3:12345", defaultLottery},
		{false, "解析应该失败-其他字符开头", "a02,03,04,05-06,07x3:12345", defaultLottery},
		{false, "解析应该失败-逗号结尾，无蓝球", "01,02,03,04,05,06x3:12345,", defaultLottery},
		{false, "解析应该失败-逗号结尾，有蓝球", "01,02,03,04,05-06x3:12345,", defaultLottery},
		{false, "解析应该失败-分割号结尾，无蓝球", "01,02,03,04,05x3:12345-", defaultLottery},
		{false, "解析应该失败-分割号结尾，有蓝球", "01,02,03,04,05-06x3:12345-", defaultLottery},
		{false, "解析应该失败-x结尾，无蓝球", "01,02,03,04,05:12345:", defaultLottery},
		{false, "解析应该失败-x结尾，有蓝球", "01,02,03,04,05-06:12345x", defaultLottery},
		{false, "解析应该失败-:结尾，无蓝球", "01,02,03,04,05x3:", defaultLottery},
		{false, "解析应该失败-:结尾，有蓝球", "01,02,03,04,05-06x3:", defaultLottery},
		{false, "解析应该失败-期号在中间", "01,02,03:12345,04,05-06x3", defaultLottery},
		{false, "解析应该失败-蓝球区重复", "01,02,03,04,05-06-07x3:12345", defaultLottery},
		{false, "解析应该失败-倍投区重复", "01,02,03,04,05-06,07x3:12345:12345", defaultLottery},
		{false, "解析应该失败-倍投区和期号区循环重复", "01,02,03,04,05-06,07x3:12345x3:12345", defaultLottery},
		{true, "解析应该成功-正确的SSQ", "01,02,03,04,05,06-07x3:12345", Lottery{"unknown", 12345, 3, []int{1, 2, 3, 4, 5, 6}, []int{7}}},
		{true, "解析应该成功-正确的DLT", "01,02,03,04,05-06,07:12345x3", Lottery{"unknown", 12345, 3, []int{1, 2, 3, 4, 5}, []int{6, 7}}},
	}

	runTests := func(tests []TestCase) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := ParseLottery(tt.input)

				if tt.isPass {
					if err != nil {
						t.Errorf("%s。输入: %s, 错误信息: %s", tt.name, tt.input, err)
					}

					if !reflect.DeepEqual(result, tt.result) {
						t.Errorf("%s: 解析结果错误。输入: %s, 预期: %+v, 实际: %+v", tt.name, tt.input, tt.result, result)
					}
				} else {
					if err == nil {
						t.Errorf("%s: 应该解析失败。输入: %s 解析结果: %v", tt.name, tt.input, result)
					}
				}
			})
		}
	}

	runTests(testsOnlyNumbers)
	runTests(testsWithScale)
	runTests(testsWithIndex)
	runTests(testsWithScaleAndIndex)
}

func TestParseNumParts(t *testing.T) {
	tests := []struct {
		shouldSuccess bool
		name          string
		input         string
		errorMsg      string
		dan           []int
		tuo           []int
	}{
		{true, "解析应该成功, 1胆1拖", "01~02", "", []int{1}, []int{2}},
		{true, "解析应该成功, 多胆1拖", "01,02~03", "", []int{1, 2}, []int{3}},
		{true, "解析应该成功, 1胆多拖", "01~02,03", "", []int{1}, []int{2, 3}},
		{true, "解析应该成功, 0胆多拖(复试)", "01,02,03", "", nil, []int{1, 2, 3}},
		{false, "解析应该失败, 有错误字符", "01,0-2,03", "号码区解析失败, 错误的字符: 【-】。", nil, nil},
		{false, "解析应该失败, 逗号开头", ",01,02,03", "号码至少为1位", nil, nil},
		{false, "解析应该失败, 逗号结尾(有胆码)", "01~02,03,", "号码至少为1位", nil, nil},
		{false, "解析应该失败, 逗号结尾(无胆码)", "01,02,03,", "号码至少为1位", nil, nil},
		{false, "解析应该失败, 波浪号开头", "~01,02,03", "号码至少为1位", nil, nil},
		{false, "解析应该失败, 波浪号结尾", "01,02,03~", "号码至少为1位", nil, nil},
		{false, "解析应该失败, 拖码区重复", "01~02,03~04", "号码区解析错误。拖区重复", nil, nil},
		{false, "解析应该失败, 复试逗号重复", "01,02,,03", "号码至少为1位", nil, nil},
		{false, "解析应该失败, 胆码区逗号重复", "01,,02~03,04", "号码至少为1位", nil, nil},
		{false, "解析应该失败, 拖码区逗号重复", "01,02~03,,04", "号码至少为1位", nil, nil},
		{false, "解析应该失败, 波浪号重复", "01,02~~03,04", "号码至少为1位", nil, nil},
		{false, "解析应该失败, 复试号码过长", "01,002,03", "号码最多为2位数", nil, nil},
		{false, "解析应该失败, 胆码区号码过长", "01,002~03", "号码最多为2位数", nil, nil},
		{false, "解析应该失败, 拖码区号码过长", "01~002,03", "号码最多为2位数", nil, nil},
		{false, "解析应该失败, 复试号码重复", "01,02,03,03,02", "号码区解析失败。拖码区重复: [2 3]。", nil, nil},
		{false, "解析应该失败, 胆码区号码重复", "01,02,02~03,04", "号码区解析失败。胆码区重复: [2]。", nil, nil},
		{false, "解析应该失败, 拖码区号码重复", "01,02~03,04,03", "号码区解析失败。拖码区重复: [3]。", nil, nil},
		{false, "解析应该失败, 胆码区号码重复, 拖码区号码重复", "01,02,01~03,04,03", "号码区解析失败。胆码区重复: [1]。拖码区重复: [3]。", nil, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dan, tuo, err := ParseNumParts(tt.input)

			if tt.shouldSuccess {
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
