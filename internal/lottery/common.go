package lottery

import (
	"errors"
	"fmt"
	"strconv"
)

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
