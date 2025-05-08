package main

import (
	"fmt"

	"github.com/buggy-95/lott/internal/lottery"
)

func main() {
	target1 := "15,18,20,21,34-04,10"
	sources := []string{
		"11,14,15,31,35-01,11",
		"02,08,09,15,30-05,08",
		"13,15,16,20,22-10,12",
		"09,15,22,23,31-07,08",
		"05,16,19,31,34-04,10",
	}

	for _, source := range sources {
		complexResult, err := lottery.GetComplexResult("DLT:"+source+"x3:25050", "DLT:"+target1)

		if err != nil {
			fmt.Printf("错误信息: %s, 输入: %s", err, source)

			continue
		}

		resultMap := map[int]string{
			1: "一等奖",
			2: "二等奖",
			3: "三等奖",
			4: "四等奖",
			5: "五等奖",
			6: "六等奖",
			7: "七等奖",
			8: "八等奖",
			9: "九等奖",
		}

		for _, result := range complexResult.List {
			var (
				front []int
				back  []int
			)

			for _, item := range result.Numbers {
				if item.Bingo {
					if item.Type == "FrontTuo" {
						front = append(front, item.Num)
					} else if item.Type == "BackTuo" {
						back = append(back, item.Num)
					}
				}
			}

			level := resultMap[result.Level]

			if len(level) == 0 {
				level = "无"
			}

			fmt.Printf("前区: %v, 后区: %v, 结果: %s\n", front, back, level)
		}
	}
}
