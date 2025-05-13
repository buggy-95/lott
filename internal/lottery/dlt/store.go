package dlt

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

func getPageData(page int) (HistoryValue, error) {
	url := fmt.Sprintf("https://webapi.sporttery.cn/gateway/lottery/getHistoryPageListV1.qry?gameNo=85&provinceId=0&pageSize=100&isVerify=1&pageNo=%d", page)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("请求失败:", err)
		return HistoryValue{}, err
	}
	defer resp.Body.Close()

	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Println("读取响应失败:", err)
	// 	return HistoryValue{}, err
	// } else {
	// 	fmt.Println("响应内容:", string(body))
	// }

	var history HistoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&history); err != nil {
		fmt.Println("解析失败, page:", page, err)
		return HistoryValue{}, err
	}

	return history.Value, nil
}

func getFullHistory(firstPageData *HistoryValue) ([]PoolDraw, error) {
	if firstPageData == nil {
		firstPage, err := getPageData(1)
		if err != nil {
			fmt.Println("历史数据获取失败:", err)
			return nil, err
		}

		firstPageData = &firstPage
	}

	concurrencyLimit := 5
	semaphore := make(chan struct{}, concurrencyLimit)
	resultChan := make(chan HistoryValue, firstPageData.Pages)

	var (
		wg     sync.WaitGroup
		mutex  sync.Mutex
		errors []error
	)

	for pageNo := 2; pageNo <= firstPageData.Pages; pageNo++ {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			fmt.Printf("total: %d, fetching: %d\n", firstPageData.Pages, pageNo)
			pageData, err := getPageData(pageNo)
			if err != nil {
				mutex.Lock()
				errors = append(errors, fmt.Errorf("页面%d请求失败: %v", p, err))
				mutex.Unlock()
				return
			}

			resultChan <- pageData
		}(pageNo)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 收集结果并按页码排序
	results := make([]HistoryValue, firstPageData.Total)
	results[0] = *firstPageData
	for resp := range resultChan {
		results[resp.PageNo-1] = resp // 页码从1开始，索引从0开始
	}

	if len(errors) > 0 {
		for _, e := range errors {
			log.Println(e)
		}

		log.Fatal("部分页面请求失败")
	}

	allData := make([]PoolDraw, firstPageData.Total)
	for _, history := range results {
		allData = append(allData, history.List...)
	}

	return allData, nil
}

func CheckStore() {
	firstPage, err := getPageData(1)
	if err != nil {
		fmt.Println("历史数据获取失败:", err)
		return
	}

	fmt.Println("数据:", firstPage.Total)
	list, listErr := getFullHistory(nil)
	if listErr != nil {
		fmt.Println("全部历史数据获取失败:", err)
		return
	}

	fullData := struct {
		UpdateTime string     `json:"updateTime"`
		List       []PoolDraw `json:"list"`
	}{
		UpdateTime: time.Now().Format("2006-01-02 15:04:05"),
		List:       list,
	}

	jsonData, jsonErr := json.MarshalIndent(fullData, "", "  ")
	if jsonErr != nil {
		fmt.Println("json解析失败:", jsonErr)
		return
	}

	writeErr := os.WriteFile("dlt_history.json", jsonData, 0644)
	if writeErr != nil {
		fmt.Println("文件写入失败:", writeErr)
		return
	}

	fmt.Println("文件写入成功")
}
