package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
)

func main() {
	client := resty.New()
	client.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	
	fmt.Println("=== 探索TDX API数据结构 ===\n")
	
	// 1. 先看看竹北站信息
	fmt.Println("1. 获取竹北站信息...")
	resp, err := client.R().
		SetQueryParam("$format", "JSON").
		SetQueryParam("$filter", "StationID eq 1180").
		Get("https://tdx.transportdata.tw/api/basic/v2/Rail/TRA/Station")
	
	if err == nil && resp.StatusCode() == 200 {
		fmt.Printf("竹北站信息: %s\n\n", resp.String())
	} else {
		fmt.Printf("获取失败: %d - %s\n\n", resp.StatusCode(), resp.String())
	}
	
	// 2. 探索GeneralTimetable数据结构
	fmt.Println("2. 探索GeneralTimetable数据结构...")
	resp, err = client.R().
		SetQueryParam("$format", "JSON").
		SetQueryParam("$top", "1").
		Get("https://tdx.transportdata.tw/api/basic/v2/Rail/TRA/GeneralTimetable")
	
	if err == nil && resp.StatusCode() == 200 {
		fmt.Printf("GeneralTimetable样例: %s\n\n", resp.String())
		
		// 尝试解析JSON结构
		var data []map[string]interface{}
		if json.Unmarshal(resp.Body(), &data) == nil && len(data) > 0 {
			fmt.Println("JSON结构分析:")
			for key, value := range data[0] {
				fmt.Printf("  %s: %T\n", key, value)
			}
		}
	} else {
		fmt.Printf("获取失败: %d - %s\n\n", resp.StatusCode(), resp.String())
	}
	
	// 3. 探索GeneralTrainInfo
	fmt.Println("\n3. 探索GeneralTrainInfo数据结构...")
	resp, err = client.R().
		SetQueryParam("$format", "JSON").
		SetQueryParam("$top", "3").
		Get("https://tdx.transportdata.tw/api/basic/v2/Rail/TRA/GeneralTrainInfo")
	
	if err == nil && resp.StatusCode() == 200 {
		fmt.Printf("GeneralTrainInfo样例: %s\n\n", resp.String())
	}
	
	// 4. 探索StationLiveBoard
	fmt.Println("4. 探索StationLiveBoard数据结构...")
	resp, err = client.R().
		SetQueryParam("$format", "JSON").
		SetQueryParam("$top", "5").
		Get("https://tdx.transportdata.tw/api/basic/v3/Rail/TRA/StationLiveBoard")
	
	if err == nil && resp.StatusCode() == 200 {
		fmt.Printf("StationLiveBoard样例: %s\n\n", resp.String())
	}
	
	// 5. 尝试过滤特定站点的实时信息
	fmt.Println("5. 尝试过滤竹北站实时信息...")
	resp, err = client.R().
		SetQueryParam("$format", "JSON").
		SetQueryParam("$filter", "StationID eq 1180").
		Get("https://tdx.transportdata.tw/api/basic/v3/Rail/TRA/StationLiveBoard")
	
	if err == nil && resp.StatusCode() == 200 {
		fmt.Printf("竹北站实时信息: %s\n", resp.String())
	} else {
		fmt.Printf("获取失败: %d - %s\n", resp.StatusCode(), resp.String())
	}
}