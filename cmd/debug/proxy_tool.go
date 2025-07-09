package main

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

func main() {
	fmt.Println("測試代理連接")
	fmt.Println("========================================")
	
	// 測試1: 不使用代理
	fmt.Println("1. 測試不使用代理:")
	testWithoutProxy()
	
	// 測試2: 使用代理
	fmt.Println("\n2. 測試使用代理:")
	testWithProxy()
	
	// 測試3: 測試代理本身
	fmt.Println("\n3. 測試代理服務器:")
	testProxyServer()
}

func testWithoutProxy() {
	client := resty.New()
	client.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	
	resp, err := client.R().
		SetQueryParam("$format", "JSON").
		SetQueryParam("$filter", "StationID eq '1180'").
		Get("https://tdx.transportdata.tw/api/basic/v3/Rail/TRA/StationLiveBoard")
	
	if err != nil {
		fmt.Printf("   ❌ 錯誤: %v\n", err)
	} else {
		fmt.Printf("   ✅ 成功: 狀態碼 %d, 大小 %d bytes\n", resp.StatusCode(), len(resp.Body()))
	}
}

func testWithProxy() {
	client := resty.New()
	client.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	client.SetProxy("http://localhost:5555/random")
	
	resp, err := client.R().
		SetQueryParam("$format", "JSON").
		SetQueryParam("$filter", "StationID eq '1180'").
		Get("https://tdx.transportdata.tw/api/basic/v3/Rail/TRA/StationLiveBoard")
	
	if err != nil {
		fmt.Printf("   ❌ 錯誤: %v\n", err)
	} else {
		fmt.Printf("   ✅ 成功: 狀態碼 %d, 大小 %d bytes\n", resp.StatusCode(), len(resp.Body()))
	}
}

func testProxyServer() {
	client := resty.New()
	
	// 直接訪問代理服務器
	resp, err := client.R().Get("http://localhost:5555/random")
	
	if err != nil {
		fmt.Printf("   ❌ 代理服務器無法訪問: %v\n", err)
	} else {
		fmt.Printf("   ✅ 代理服務器回應: 狀態碼 %d\n", resp.StatusCode())
		fmt.Printf("   響應內容: %s\n", resp.String())
	}
	
	// 測試代理根路徑
	resp2, err2 := client.R().Get("http://localhost:5555")
	
	if err2 != nil {
		fmt.Printf("   ❌ 代理根路徑無法訪問: %v\n", err2)
	} else {
		fmt.Printf("   ✅ 代理根路徑回應: 狀態碼 %d\n", resp2.StatusCode())
		fmt.Printf("   響應內容: %s\n", resp2.String())
	}
}