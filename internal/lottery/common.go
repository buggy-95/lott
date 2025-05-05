package lottery

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
)

func getDupNums(nums []int) []int {
	var dupSlice []int

	dupMap := make(map[int]bool)

	for _, num := range nums {
		if dupMap[num] {
			dupSlice = append(dupSlice, num)
		} else {
			dupMap[num] = true
		}
	}

	sort.Ints(dupSlice)

	return dupSlice
}

func ParseLottery(input string) (Lottery, error) {
	nextInputType := "red" // 下一个数字应该是什么类型的，流转顺序 red -> blue -> scale | index
	red := []int{}
	blue := []int{}
	numStr := ""
	scaleParsed := false
	indexParsed := false
	errorMsg := fmt.Sprintf("Invalid lottery: %s", input)

	lottery := Lottery{
		Type:  "unknown",
		Index: 0,
		Red:   []int{},
		Blue:  []int{},
		Scale: 1,
	}

	switchNextInputType := func(next string) bool {
		switch next {
		case "blue":
			if nextInputType != "red" {
				return false
			}
		case "scale":
			if scaleParsed || !(nextInputType == "blue" || nextInputType == "index") {
				return false
			} else {
				scaleParsed = true
			}
		case "index":
			if indexParsed || !(nextInputType == "blue" || nextInputType == "scale") {
				return false
			} else {
				indexParsed = false
			}
		}

		nextInputType = next

		return true
	}

	dealNumber := func() bool {
		if len(numStr) < 1 {
			return false
		}

		num, _ := strconv.Atoi(numStr)

		switch nextInputType {
		case "red":
			red = append(red, num)
		case "blue":
			blue = append(blue, num)
		case "scale":
			lottery.Scale = num
		case "index":
			lottery.Index = num
		default:
			return false
		}

		numStr = ""

		return true
	}

	for _, char := range input {
		if '0' <= char && char <= '9' {
			numStr += string(char)

			// 每个数字最多2位
			if (nextInputType == "red" || nextInputType == "blue") && len(numStr) > 2 {
				return lottery, errors.New(errorMsg)
			}

			continue
		}

		if succ := dealNumber(); !succ {
			return lottery, errors.New(errorMsg)
		}

		if char == ',' {
			continue
		}

		nextTypeMap := map[rune]string{
			'-': "blue",
			'x': "scale",
			':': "index",
		}

		if succ := switchNextInputType(nextTypeMap[char]); succ {
			continue
		}

		return lottery, errors.New(errorMsg)
	}

	if succ := dealNumber(); !succ || nextInputType == "red" {
		return lottery, errors.New(errorMsg)
	}

	lottery.Red = red
	lottery.Blue = blue

	return lottery, nil
}

// 将号码区解析为前区和后区
func ParseNumParts(input string) ([]int, []int, error) {
	var (
		dan []int  // 胆码
		tuo []int  // 拖码
		str string // 缓冲区
	)

	nextNumberType := "dan" // dan -> tuo

	switchNextNumberType := func(next string) error {
		if next == "tuo" {
			if nextNumberType != "dan" {
				return fmt.Errorf("号码区解析错误。拖区重复。输入: %s", input)
			}

			nextNumberType = "tuo"

			return nil
		}

		return fmt.Errorf("号码区解析错误, 不支持的状态流转: %s。输入: %s", next, input)
	}

	dealNumber := func() error {
		if len(str) < 1 {
			return fmt.Errorf("号码至少为1位")
		}

		num, err := strconv.Atoi(str)

		if err != nil {
			return fmt.Errorf("号码解析失败，错误信息: %s", err)
		}

		switch nextNumberType {
		case "dan":
			dan = append(dan, num)
		case "tuo":
			tuo = append(tuo, num)
		default:
			return fmt.Errorf("号码解析失败，错误的类型: %s", nextNumberType)
		}

		// 号码处理完成后清除缓冲区
		str = ""

		return nil
	}

	for _, char := range input {
		if '0' <= char && char <= '9' {
			// 追加数字字符
			str += string(char)

			// 缓冲字符长度超过2抛错
			if len(str) > 2 {
				return nil, nil, errors.New("号码最多为2位数")
			}

			continue
		}

		// 遇到非数字的字符尝试处理缓冲区的数字
		if err := dealNumber(); err != nil {
			return nil, nil, err
		}

		if char == ',' {
			continue
		}

		if char == '~' {
			err := switchNextNumberType("tuo")

			if err != nil {
				return nil, nil, err
			}

			continue
		}

		return nil, nil, fmt.Errorf("号码区解析失败, 错误的字符: 【%c】。输入: %s", char, input)
	}

	if err := dealNumber(); err != nil {
		return nil, nil, err
	}

	if len(dan) == 0 {
		return nil, nil, fmt.Errorf("号码区解析失败。输入: %s", input)
	}

	if len(tuo) == 0 {
		tuo, dan = dan, tuo
	}

	dupDan := getDupNums(dan)
	dupTuo := getDupNums(tuo)
	dupMsg := ""

	if len(dupDan) > 0 {
		dupMsg += fmt.Sprintf("胆码区重复: %v。", dupDan)
	}

	if len(dupTuo) > 0 {
		dupMsg += fmt.Sprintf("拖码区重复: %v。", dupTuo)
	}

	if len(dupMsg) > 0 {
		dupMsg = "号码区解析失败。" + dupMsg

		return nil, nil, errors.New(dupMsg)
	}

	return dan, tuo, nil
}

// DLT:01,02,03,04~05,06-07~08x3:25053
func ParseComplexLotteryParts(input string) (ComplexLotteryParts, error) {
	complexLotteryParts := ComplexLotteryParts{}

	return complexLotteryParts, nil
}
