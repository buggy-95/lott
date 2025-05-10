package lottery

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
)

// GetDupNums
//
// @Description 获取重复的号码
//
// @Param nums []int 数字列表
//
// @Return []int 重复的号码列表
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

// GetCrossNums
//
// @Description 获取两个数字列表的交集
//
// @Param source []int 数字列表1
//
// @Param target []int 数字列表2
//
// @Return []int 交集列表
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

// parseNumParts
//
// @Description 将数字区解析为胆码区和拖码区，通过波浪号将数字区进行分隔，若无波浪号则认为整个数字区都是拖码区。
//
// @Param input string 输入的号码区字符串，格式为: 01,02～03,04 或 01,02,03
//
// @Return []int 胆码区的数字列表
//
// @Return []int 拖码区的数字列表
//
// @Return error 错误信息
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

// parseLotteryParts
//
// @Description 解析复杂彩票的字符串，格式为: 彩票类型: 前区号码-后区号码[x倍投][:期号]
//
// @Param input string 输入的复杂彩票字符串，例如：DLT:01,02,03,04~05,06-07~08x3:25053
//
// @Return ComplexLotteryParts 解析后的复杂彩票结构体
//
// @Return error 错误信息
func parseLotteryParts(input string) (LotteryParts, error) {
	nextTokenType := "type" // type -> front -> back -> scale | index -> index | scale
	lotteryParts := LotteryParts{}
	lotteryParts.Scale = 1

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
				lotteryParts.Type = tmpToken
			} else {
				return fmt.Errorf("彩票类型解析失败。不支持的彩票类型: %s。输入: %s", tmpToken, input)
			}
		case "front":
			dan, tuo, err := parseNumParts(tmpToken)

			if err != nil {
				return fmt.Errorf("前区号码解析失败。原因: %s。输入: %s", err, input)
			} else {
				lotteryParts.FrontDan = dan
				lotteryParts.FrontTuo = tuo
			}
		case "back":
			dan, tuo, err := parseNumParts(tmpToken)

			if err != nil {
				return fmt.Errorf("后区号码解析失败。原因: %s。输入: %s", err, input)
			} else {
				lotteryParts.BackDan = dan
				lotteryParts.BackTuo = tuo
			}
		case "scale":
			scale, err := strconv.Atoi(tmpToken)

			if err != nil {
				return fmt.Errorf("倍投倍数解析失败。倍数: %s。输入: %s", tmpToken, input)
			} else {
				lotteryParts.Scale = scale
				scaleParsed = true
			}
		case "index":
			index, err := strconv.Atoi(tmpToken)

			if err != nil {
				return fmt.Errorf("期号解析失败。期号: %s。输入: %s", tmpToken, input)
			} else {
				lotteryParts.Index = index
				indexParsed = true
			}
		}

		return nil
	}

	dealChar := func(char rune) error {
		isDigit := func(char rune) bool {
			return '0' <= char && char <= '9'
		}

		isUpperAlpha := func(char rune) bool {
			return 'A' <= char && char <= 'Z'
		}

		appendToken := func(char rune) {
			token += string(char)
		}

		handleTransition := func(expectedTokenType string) error {
			if err := dealToken(nextTokenType); err != nil {
				return err
			}

			return switchNextTokenType(expectedTokenType)
		}

		switch nextTokenType {
		case "type":
			if char == ':' {
				return handleTransition("front")
			} else if len(token) < 5 && isUpperAlpha(char) {
				appendToken(char)
			} else {
				return fmt.Errorf("彩票类型解析失败。输入: %s", input)
			}
		case "front":
			if char == '-' {
				return handleTransition("back")
			} else if char == ',' || char == '~' || isDigit(char) {
				appendToken(char)
			} else {
				return fmt.Errorf("前区号码解析失败。当前字符: 【%c】。输入: %s", char, input)
			}
		case "back":
			if char == 'x' {
				return handleTransition("scale")
			} else if char == ':' {
				return handleTransition("index")
			} else if char == ',' || char == '~' || isDigit(char) {
				appendToken(char)
			} else {
				return fmt.Errorf("后区号码解析失败。输入: %s", input)
			}
		case "scale":
			if char == ':' {
				return handleTransition("index")
			} else if isDigit(char) {
				appendToken(char)
			} else {
				return fmt.Errorf("倍投数解析失败。当前字符: 【%c】。输入: %s", char, input)
			}
		case "index":
			if char == 'x' {
				return handleTransition("scale")
			} else if isDigit(char) {
				appendToken(char)
			} else {
				return fmt.Errorf("期数解析失败。当前字符: 【%c】。输入: %s", char, input)
			}
		}

		return nil
	}

	for _, char := range input {
		if err := dealChar(char); err != nil {
			return LotteryParts{}, err
		}
	}

	if err := dealToken(nextTokenType); err != nil {
		return LotteryParts{}, err
	}

	return lotteryParts, nil
}

// genPermutation
//
// @Description 从列表中生成长度为n的所有组合，按照从小到大的顺序排列
//
// @Param nums []int 数字列表
//
// @Param n int 组合的长度
//
// @Return [][]int 组合列表
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

// genLotteryList
//
// @Description 生成单式彩票列表
//
// @Return []SingleLottery 单式彩票列表
func (parts *LotteryParts) genLotteryList() []Lottery {
	var (
		result    []Lottery
		frontList [][]int
		backList  [][]int
	)

	baseInfo := parts.LotteryBaseInfo
	frontDan := parts.FrontDan
	frontTuo := parts.FrontTuo
	backDan := parts.BackDan
	backTuo := parts.BackTuo

	// TODO: 考虑这里的逻辑是否可以拆到外面去
	switch baseInfo.Type {
	case "DLT":
		frontList = genPermutation(frontTuo, 5-len(frontDan))
		backList = genPermutation(backTuo, 2-len(backDan))
	case "SSQ":
		frontList = genPermutation(frontTuo, 6-len(frontDan))
		backList = genPermutation(backTuo, 1-len(backDan))
	}

	for _, front := range frontList {
		for _, back := range backList {
			lott := Lottery{
				LotteryParts: LotteryParts{
					LotteryBaseInfo: baseInfo,
					FrontTuo:        append(frontDan, front...),
					BackTuo:         append(backDan, back...),
					FrontDan:        nil,
					BackDan:         nil,
				},
			}

			sort.Ints(lott.FrontTuo)
			sort.Ints(lott.BackTuo)

			result = append(result, lott)
		}
	}

	return result
}

// getMatchNums
//
// @Description 获取两个数字列表的交集，返回标记是否命中的source号码列表和命中数量
//
// @Param source []int source数字列表
//
// @Param target []int target数字列表
//
// @Return []BingoNum 命中号码列表
//
// @Return int 命中数量
func getMatchNums(source []int, target []int) ([]BingoNum, int) {
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

// isSingleLottery
//
// @Description 判断当前彩票是否是单式票
//
// @Return bool 当前彩票是否是单式票
func (lott *Lottery) IsSingleLottery() bool {
	return len(lott.List) == 0
}

// getLotteryResult
//
// @Description 获取彩票的开奖结果
//
// @Param target Lottery 开奖彩票，必须要是单式票
//
// @Return LotteryResult 购奖彩票的开奖结果
//
// @Return error 错误信息
func (source *Lottery) GetLotteryResult(target Lottery) (LotteryResult, error) {
	var (
		result LotteryResult
		nums   []ResultNum
		level  int
	)

	if !target.IsSingleLottery() {
		return result, fmt.Errorf("开奖彩票不是单式票: %s", target.format(true))
	}

	result.LotteryBaseInfo = source.LotteryBaseInfo

	// 处理单式票结果
	if source.IsSingleLottery() {
		frontNums, frontMatched := getMatchNums(source.FrontTuo, target.FrontTuo)
		backNums, backMatched := getMatchNums(source.BackTuo, target.BackTuo)

		for _, num := range frontNums {
			nums = append(nums, ResultNum{Type: "FrontTuo", BingoNum: num})
		}

		for _, num := range backNums {
			nums = append(nums, ResultNum{Type: "BackTuo", BingoNum: num})
		}

		switch frontMatched {
		case 5:
			switch backMatched {
			case 2:
				level = 1
			case 1:
				level = 2
			case 0:
				level = 3
			}
		case 4:
			switch backMatched {
			case 2:
				level = 4
			case 1:
				level = 5
			case 0:
				level = 7
			}
		case 3:
			switch backMatched {
			case 2:
				level = 6
			case 1:
				level = 8
			case 0:
				level = 9
			}
		case 2:
			switch backMatched {
			case 2:
				level = 8
			case 1:
				level = 9
			}
		case 1:
			switch backMatched {
			case 2:
				level = 9
			}
		case 0:
			switch backMatched {
			case 2:
				level = 9
			}
		}

		levelPriceMap := map[int]int{
			1: 10000000,
			2: 200000,
			3: 10000,
			4: 3000,
			5: 300,
			6: 200,
			7: 100,
			8: 15,
			9: 5,
		}

		result.LotteryBaseInfo = source.LotteryBaseInfo
		result.FrontMatched = frontMatched
		result.BackMatched = backMatched
		result.Numbers = nums
		result.Level = level
		result.Price = levelPriceMap[level] * source.Scale

		return result, nil
	}

	// 复式票递归处理
	frontDanNums, frontDanMatched := getMatchNums(source.FrontDan, target.FrontTuo)
	frontTuoNums, frontTuoMatched := getMatchNums(source.FrontTuo, target.FrontTuo)
	backDanNums, backDanMatched := getMatchNums(source.BackDan, target.BackTuo)
	backTuoNums, backTuoMatched := getMatchNums(source.BackTuo, target.BackTuo)

	result.FrontMatched = frontDanMatched + frontTuoMatched
	result.BackMatched = backDanMatched + backTuoMatched

	for _, num := range frontDanNums {
		nums = append(nums, ResultNum{Type: "FrontDan", BingoNum: num})
	}

	for _, num := range frontTuoNums {
		nums = append(nums, ResultNum{Type: "FrontTuo", BingoNum: num})
	}

	for _, num := range backDanNums {
		nums = append(nums, ResultNum{Type: "BackDan", BingoNum: num})
	}

	for _, num := range backTuoNums {
		nums = append(nums, ResultNum{Type: "BackTuo", BingoNum: num})
	}

	result.Numbers = nums
	result.Level = 100

	for _, lott := range source.List {
		lottResult, err := lott.GetLotteryResult(target)

		if err != nil {
			fmt.Println(err)
			return result, err
		}

		result.List = append(result.List, lottResult)
		result.Price += lottResult.Price

		if lottResult.Level > 0 {
			result.Level = min(result.Level, lottResult.Level)
		}
	}

	if result.Level == 100 {
		result.Level = 0
	}

	return result, nil
}

// GetLottery
//
// @Description 获取复杂彩票的结构体
//
// @Param input string 输入的复杂彩票字符串，例如：DLT:01,02,03,04~05,06-07~08x3:25053
//
// @Return Lottery 复杂彩票结构体
//
// @Return error 错误信息
func GetLottery(input string) (Lottery, error) {
	var (
		result Lottery
	)

	parts, err := parseLotteryParts(input)

	if err != nil {
		return result, err
	}

	list := parts.genLotteryList()

	if len(list) > 1 {
		return Lottery{parts, list}, nil
	} else if len(list) == 1 {
		return Lottery{parts, nil}, nil
	} else {
		return result, fmt.Errorf("单式票生成失败，输入: %s", input)
	}
}

// format
//
// @Description 格式化彩票的号码区，例如: 01,02,03,04-05,06
//
// @Param showExtra bool 是否展示除了号码区以外的额外信息，例如倍投倍数和期号
//
// @Return string 格式化后的号码区字符串
func (lott *Lottery) format(showExtra bool) string {
	var (
		front string
		back  string
		str   string
	)

	// 前区胆码
	for _, num := range lott.FrontDan {
		front += fmt.Sprintf(",%02d", num)
	}

	if len(front) > 0 {
		str += front[1:]
		front = ""
	}

	// 前区拖码
	for _, num := range lott.FrontTuo {
		front += fmt.Sprintf(",%02d", num)
	}

	if len(front) > 0 {
		if len(lott.FrontDan) > 0 {
			str += "~"
		}

		str += front[1:] + "-"
	}

	// 后区胆码
	for _, num := range lott.BackDan {
		back += fmt.Sprintf(",%02d", num)
	}

	if len(back) > 0 {
		str += back[1:]
		back = ""
	}

	// 后区拖码
	for _, num := range lott.BackTuo {
		back += fmt.Sprintf(",%02d", num)
	}

	if len(back) > 0 {
		if len(lott.BackDan) > 0 {
			str += "~"
		}

		str += back[1:]
	}

	if !showExtra {
		return str
	}

	if scale := lott.LotteryBaseInfo.Scale; scale > 1 {
		str += fmt.Sprintf("x%d", scale)
	}

	if index := lott.LotteryBaseInfo.Index; index > 0 {
		str += fmt.Sprintf(":%d", index)
	}

	return str
}
