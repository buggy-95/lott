package lottery

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
)

func GetDupNums(nums []int) []int {
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

func GetCrossNums(source, target []int) []int {
	var result []int

	numMap := make(map[int]bool, len(source))

	for _, num := range source {
		numMap[num] = true
	}

	for _, num := range target {
		if numMap[num] {
			result = append(result, num)
		}
	}

	sort.Ints(result)

	return result
}

// 将号码区解析为前区和后区
func parseNumParts(input string) ([]int, []int, error) {
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
		} else {
			return fmt.Errorf("号码区解析错误，不支持的状态流转: %s。输入: %s", next, input)
		}

		return nil
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

	check := func() error {
		if len(dan) == 0 {
			return fmt.Errorf("号码区解析失败。输入: %s", input)
		}

		if len(tuo) == 0 {
			tuo, dan = dan, tuo
		}

		dupMsg := ""

		if dupDan := GetDupNums(dan); len(dupDan) > 0 {
			dupMsg += fmt.Sprintf("胆码区重复: %v。", dupDan)
		}

		if dupTuo := GetDupNums(tuo); len(dupTuo) > 0 {
			dupMsg += fmt.Sprintf("拖码区重复: %v。", dupTuo)
		}

		if len(dupMsg) > 0 {
			dupMsg = "号码区解析失败。" + dupMsg

			return errors.New(dupMsg)
		}

		if cross := GetCrossNums(dan, tuo); len(cross) > 0 {
			return fmt.Errorf("号码区解析失败。拖码区冲突: %v。", cross)
		}

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

		return nil, nil, fmt.Errorf("号码区解析失败，错误的字符: 【%c】。输入: %s", char, input)
	}

	if err := dealNumber(); err != nil {
		return nil, nil, err
	}

	if err := check(); err != nil {
		return nil, nil, err
	}

	return dan, tuo, nil
}

// DLT:01,02,03,04~05,06-07~08x3:25053
func ParseComplexLotteryParts(input string) (ComplexLotteryParts, error) {
	nextTokenType := "type" // type -> front -> back -> scale | index -> index | scale
	complexLotteryParts := ComplexLotteryParts{}
	complexLotteryParts.Scale = 1

	var (
		token       string
		scaleParsed bool
		indexParsed bool
	)

	switchNextTokenType := func(next string) error {
		commonErrorMsg := fmt.Errorf("解析失败。错误的token类型: %s, 当前token类型: %s。输入: %s", next, nextTokenType, input)

		switch next {
		case "front":
			if nextTokenType != "type" {
				return commonErrorMsg
			}
		case "back":
			if nextTokenType != "front" {
				return commonErrorMsg
			}
		case "scale":
			if scaleParsed {
				return fmt.Errorf("解析失败。倍投已解析过。输入: %s", input)
			}

			if !(nextTokenType == "back" || nextTokenType == "index") {
				return commonErrorMsg
			}
		case "index":
			if indexParsed {
				return fmt.Errorf("解析失败。期号已解析过。输入: %s", input)
			}

			if !(nextTokenType == "back" || nextTokenType == "scale") {
				return commonErrorMsg
			}
		}

		nextTokenType = next

		return nil
	}

	dealToken := func(tokenType string) error {
		tmpToken := token
		token = ""

		switch tokenType {
		case "type":
			if tmpToken == "SSQ" || tmpToken == "DLT" {
				complexLotteryParts.Type = tmpToken
			} else {
				return fmt.Errorf("彩票类型解析失败。不支持的彩票类型: %s。输入: %s", tmpToken, input)
			}
		case "front":
			dan, tuo, err := parseNumParts(tmpToken)

			if err != nil {
				return fmt.Errorf("前区号码解析失败。原因: %s。输入: %s", err, input)
			} else {
				complexLotteryParts.FrontDan = dan
				complexLotteryParts.FrontTuo = tuo
			}
		case "back":
			dan, tuo, err := parseNumParts(tmpToken)

			if err != nil {
				return fmt.Errorf("后区号码解析失败。原因: %s。输入: %s", err, input)
			} else {
				complexLotteryParts.BackDan = dan
				complexLotteryParts.BackTuo = tuo
			}
		case "scale":
			scale, err := strconv.Atoi(tmpToken)

			if err != nil {
				return fmt.Errorf("倍投倍数解析失败。倍数: %s。输入: %s", tmpToken, input)
			} else {
				complexLotteryParts.Scale = scale
				scaleParsed = true
			}
		case "index":
			index, err := strconv.Atoi(tmpToken)

			if err != nil {
				return fmt.Errorf("期号解析失败。期号: %s。输入: %s", tmpToken, input)
			} else {
				complexLotteryParts.Index = index
				indexParsed = true
			}
		}

		return nil
	}

	dealChar := func(char rune) error {
		switch nextTokenType {
		case "type":
			if char == ':' {
				if err := dealToken("type"); err != nil {
					return err
				}

				if err := switchNextTokenType("front"); err != nil {
					return err
				}
			} else if len(token) < 5 && ('A' <= char && char <= 'Z') {
				token += string(char)
			} else {
				return fmt.Errorf("彩票类型解析失败。输入: %s", input)
			}
		case "front":
			if char == '-' {
				if err := dealToken("front"); err != nil {
					return err
				}

				if err := switchNextTokenType("back"); err != nil {
					return err
				}
			} else if char == ',' || char == '~' || ('0' <= char && char <= '9') {
				token += string(char)
			} else {
				return fmt.Errorf("前区号码解析失败。当前字符: 【%c】。输入: %s", char, input)
			}
		case "back":
			if char == 'x' {
				if err := dealToken("back"); err != nil {
					return err
				}

				if err := switchNextTokenType("scale"); err != nil {
					return err
				}
			} else if char == ':' {
				if err := dealToken("back"); err != nil {
					return err
				}

				if err := switchNextTokenType("index"); err != nil {
					return err
				}
			} else if char == ',' || char == '~' || ('0' <= char && char <= '9') {
				token += string(char)
			} else {
				return fmt.Errorf("后区号码解析失败。输入: %s", input)
			}
		case "scale":
			if char == ':' {
				if err := dealToken("scale"); err != nil {
					return err
				}

				if err := switchNextTokenType("index"); err != nil {
					return err
				}
			} else if '0' <= char && char <= '9' {
				token += string(char)
			} else {
				return fmt.Errorf("倍投数解析失败。当前字符: 【%c】。输入: %s", char, input)
			}
		case "index":
			if char == 'x' {
				if err := dealToken("index"); err != nil {
					return err
				}

				if err := switchNextTokenType("scale"); err != nil {
					return err
				}
			} else if '0' <= char && char <= '9' {
				token += string(char)
			} else {
				return fmt.Errorf("期数解析失败。当前字符: 【%c】。输入: %s", char, input)
			}
		}

		return nil
	}

	for _, char := range input {
		if err := dealChar(char); err != nil {
			return ComplexLotteryParts{}, err
		}
	}

	if err := dealToken(nextTokenType); err != nil {
		return ComplexLotteryParts{}, err
	}

	return complexLotteryParts, nil
}

func genPermutation(nums []int, n int) [][]int {
	var (
		result    [][]int
		backtrack func(start int, current []int)
	)

	if n < 0 || n >= len(nums) {
		return [][]int{nums}
	}

	backtrack = func(start int, current []int) {
		if len(current) == n {
			temp := make([]int, n)
			copy(temp, current)
			result = append(result, temp)

			return
		}

		for i := start; i < len(nums); i++ {
			current = append(current, nums[i])
			backtrack(i+1, current)
			current = current[:len(current)-1]
		}
	}

	backtrack(0, []int{})

	return result
}

func genSingleLotteryList(parts ComplexLotteryParts) []SingleLottery {
	var (
		result []SingleLottery
	)

	frontDan := parts.FrontDan
	frontTuo := parts.FrontTuo
	backDan := parts.BackDan
	backTuo := parts.BackTuo

	frontList := genPermutation(frontTuo, 5-len(frontDan))
	backList := genPermutation(backTuo, 2-len(backDan))

	for _, front := range frontList {
		for _, back := range backList {
			singleLottery := SingleLottery{
				LotteryBaseInfo: LotteryBaseInfo{
					Type:  parts.Type,
					Index: parts.Index,
					Scale: parts.Scale,
				},
				Front: append(frontDan, front...),
				Back:  append(backDan, back...),
			}

			sort.Ints(singleLottery.Front)
			sort.Ints(singleLottery.Back)

			result = append(result, singleLottery)
		}
	}

	return result
}

func getMatchResult(source []int, target []int) ([]BingoNum, int) {
	var (
		result  []BingoNum
		matched int
	)

	numMap := make(map[int]bool, len(target))

	for _, num := range target {
		numMap[num] = true
	}

	for _, num := range source {
		result = append(result, BingoNum{Num: num, Bingo: numMap[num]})

		if numMap[num] {
			matched++
		}
	}

	return result, matched
}
