package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
)

func main() {
	client := resty.New()
	client.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	
	baseURL := "https://tdx.transportdata.tw/api/basic/v3"
	endpoint := "/Rail/TRA/StationLiveBoard"
	
	// 創建響應日誌文件
	logFile, err := os.Create("responses.log")
	if err != nil {
		fmt.Printf("無法創建日誌文件: %v\n", err)
		return
	}
	defer logFile.Close()
	
	fmt.Println("測試 TDX API 請求限制")
	fmt.Println("API: " + baseURL + endpoint)
	fmt.Println("測試參數: StationID eq '1180'")
	fmt.Println("響應日誌: responses.log")
	fmt.Println("========================================")
	
	successCount := 0
	errorCount := 0
	rateLimitCount := 0
	
	for i := 1; i <= 200; i++ {
		startTime := time.Now()
		
		resp, err := client.R().
			SetQueryParam("$format", "JSON").
			SetQueryParam("$filter", "StationID eq '1180'").
			Get(baseURL + endpoint)
		
		duration := time.Since(startTime)
		
		// 記錄響應到文件
		logFile.WriteString(fmt.Sprintf("=== 請求 %d ===\n", i))
		logFile.WriteString(fmt.Sprintf("時間: %s\n", time.Now().Format("2006-01-02 15:04:05")))
		logFile.WriteString(fmt.Sprintf("耗時: %v\n", duration))
		
		if err != nil {
			errorCount++
			logFile.WriteString(fmt.Sprintf("錯誤: %v\n", err))
			logFile.WriteString("\n")
			fmt.Printf("請求 %d: ❌ 網路錯誤 (耗時: %v) - %v\n", i, duration, err)
		} else {
			statusCode := resp.StatusCode()
			
			// 記錄響應詳情
			logFile.WriteString(fmt.Sprintf("狀態碼: %d\n", statusCode))
			logFile.WriteString(fmt.Sprintf("響應大小: %d bytes\n", len(resp.Body())))
			
			// 記錄響應頭
			logFile.WriteString("響應頭:\n")
			for key, values := range resp.Header() {
				for _, value := range values {
					logFile.WriteString(fmt.Sprintf("  %s: %s\n", key, value))
				}
			}
			
			// 記錄響應內容（前500字符）
			responseBody := resp.String()
			if len(responseBody) > 500 {
				logFile.WriteString(fmt.Sprintf("響應內容 (前500字符): %s...\n", responseBody[:500]))
			} else {
				logFile.WriteString(fmt.Sprintf("響應內容: %s\n", responseBody))
			}
			logFile.WriteString("\n")
			
			if statusCode == 200 {
				successCount++
				fmt.Printf("請求 %d: ✅ 成功 (狀態: %d, 耗時: %v, 大小: %d bytes)\n", 
					i, statusCode, duration, len(resp.Body()))
			} else if statusCode == 429 {
				rateLimitCount++
				fmt.Printf("請求 %d: 🚫 達到限制! (狀態: %d, 耗時: %v)\n", i, statusCode, duration)
				fmt.Printf("響應內容: %s\n", resp.String())
				break
			} else {
				errorCount++
				fmt.Printf("請求 %d: ❌ HTTP錯誤 (狀態: %d, 耗時: %v)\n", i, statusCode, duration)
				fmt.Printf("響應內容: %s\n", resp.String())
			}
		}
		
		// 每 10 次請求顯示統計
		if i%10 == 0 {
			fmt.Printf("--- 進度: %d/200, 成功: %d, 錯誤: %d, 限制: %d ---\n", 
				i, successCount, errorCount, rateLimitCount)
		}
		
		// 如果已經觸發限制，就停止
		if rateLimitCount > 0 {
			break
		}
		
		// 短暫延遲
		time.Sleep(50 * time.Millisecond)
	}
	
	fmt.Println("\n========================================")
	fmt.Println("最終統計:")
	fmt.Printf("成功請求: %d\n", successCount)
	fmt.Printf("一般錯誤: %d\n", errorCount)
	fmt.Printf("限制錯誤: %d\n", rateLimitCount)
	fmt.Printf("總請求數: %d\n", successCount+errorCount+rateLimitCount)
	
	if rateLimitCount > 0 {
		fmt.Printf("🎯 在第 %d 次請求時觸發了限制！\n", successCount+errorCount+rateLimitCount)
	} else {
		fmt.Println("🤔 在 200 次請求內沒有觸發限制")
	}
	
	fmt.Println("詳細響應已保存到 responses.log")
}