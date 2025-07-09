package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"tg-rail-shouting/internal/config"
	"tg-rail-shouting/internal/tdx"
)

func main() {
	logrus.SetLevel(logrus.InfoLevel)
	
	// 載入配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	
	// 創建 TDX 客戶端
	tdxClient := tdx.NewClient(
		cfg.TDX.ClientID,
		cfg.TDX.ClientSecret,
		cfg.TDX.BaseURL,
		cfg.TDX.AuthURL,
	)
	
	fmt.Println("開始測試 API 請求限制...")
	fmt.Println("目標 API: TDX StationLiveBoard")
	fmt.Println("測試參數: 竹北站 (ID: 1180)")
	fmt.Println("========================================")
	
	// 測試單線程連續請求
	fmt.Println("\n1. 測試單線程連續請求:")
	testSequentialRequests(tdxClient, cfg.Station.ZhubeiStationID)
	
	// 等待一段時間
	fmt.Println("\n等待 30 秒後進行並發測試...")
	time.Sleep(30 * time.Second)
	
	// 測試並發請求
	fmt.Println("\n2. 測試並發請求:")
	testConcurrentRequests(tdxClient, cfg.Station.ZhubeiStationID)
	
	fmt.Println("\n測試完成！")
}

func testSequentialRequests(client *tdx.Client, stationID string) {
	successCount := 0
	errorCount := 0
	rateLimitCount := 0
	
	fmt.Printf("連續請求測試 (每次間隔 100ms):\n")
	
	for i := 1; i <= 100; i++ {
		startTime := time.Now()
		
		_, err := client.GetTrainTimetable(stationID, 1)
		duration := time.Since(startTime)
		
		if err != nil {
			errorCount++
			if contains(err.Error(), "429") || contains(err.Error(), "rate limit") {
				rateLimitCount++
				fmt.Printf("請求 %d: ❌ 觸發限制! (耗時: %v) - %v\n", i, duration, err)
				break
			} else {
				fmt.Printf("請求 %d: ❌ 錯誤 (耗時: %v) - %v\n", i, duration, err)
			}
		} else {
			successCount++
			fmt.Printf("請求 %d: ✅ 成功 (耗時: %v)\n", i, duration)
		}
		
		time.Sleep(100 * time.Millisecond)
	}
	
	fmt.Printf("\n結果統計:\n")
	fmt.Printf("成功請求: %d\n", successCount)
	fmt.Printf("一般錯誤: %d\n", errorCount-rateLimitCount)
	fmt.Printf("限制錯誤: %d\n", rateLimitCount)
}

func testConcurrentRequests(client *tdx.Client, stationID string) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	
	successCount := 0
	errorCount := 0
	rateLimitCount := 0
	
	concurrency := 10
	requestsPerGoroutine := 5
	
	fmt.Printf("並發請求測試 (%d 個並發, 每個發送 %d 個請求):\n", concurrency, requestsPerGoroutine)
	
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			
			for j := 0; j < requestsPerGoroutine; j++ {
				startTime := time.Now()
				
				_, err := client.GetTrainTimetable(stationID, 1)
				duration := time.Since(startTime)
				
				mu.Lock()
				requestNum := goroutineID*requestsPerGoroutine + j + 1
				
				if err != nil {
					errorCount++
					if contains(err.Error(), "429") || contains(err.Error(), "rate limit") {
						rateLimitCount++
						fmt.Printf("協程 %d 請求 %d: ❌ 觸發限制! (耗時: %v) - %v\n", goroutineID, requestNum, duration, err)
					} else {
						fmt.Printf("協程 %d 請求 %d: ❌ 錯誤 (耗時: %v) - %v\n", goroutineID, requestNum, duration, err)
					}
				} else {
					successCount++
					fmt.Printf("協程 %d 請求 %d: ✅ 成功 (耗時: %v)\n", goroutineID, requestNum, duration)
				}
				mu.Unlock()
				
				time.Sleep(50 * time.Millisecond)
			}
		}(i)
	}
	
	wg.Wait()
	
	fmt.Printf("\n並發測試結果統計:\n")
	fmt.Printf("成功請求: %d\n", successCount)
	fmt.Printf("一般錯誤: %d\n", errorCount-rateLimitCount)
	fmt.Printf("限制錯誤: %d\n", rateLimitCount)
}

func contains(str, substr string) bool {
	return len(str) >= len(substr) && 
		   (str == substr || 
		    str[:len(substr)] == substr || 
		    str[len(str)-len(substr):] == substr ||
		    findSubstring(str, substr))
}

func findSubstring(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}