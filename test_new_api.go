package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

func testEndpoint(client *resty.Client, name, url string) {
	fmt.Printf("\n=== 测试 %s ===\n", name)
	fmt.Printf("URL: %s\n", url)
	
	resp, err := client.R().
		SetQueryParam("$format", "JSON").
		SetQueryParam("$top", "3").
		Get(url)
	
	fmt.Printf("状态码: %d\n", resp.StatusCode())
	if err != nil {
		fmt.Printf("错误: %v\n", err)
	} else if resp.StatusCode() == 200 {
		fmt.Printf("成功! 数据长度: %d 字符\n", len(resp.String()))
		fmt.Printf("前100字符: %s\n", resp.String()[:min(100, len(resp.String()))])
	} else {
		fmt.Printf("失败! 响应: %s\n", resp.String())
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	client := resty.New()
	client.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	
	fmt.Println("TDX API 端点测试")
	
	// 测试站点信息
	testEndpoint(client, "站点信息 (v2)", "https://tdx.transportdata.tw/api/basic/v2/Rail/TRA/Station")
	testEndpoint(client, "站点信息 (v3)", "https://tdx.transportdata.tw/api/basic/v3/Rail/TRA/Station")
	
	// 测试竹北站特定信息
	testEndpoint(client, "竹北站信息", "https://tdx.transportdata.tw/api/basic/v2/Rail/TRA/Station?$filter=StationID eq '1180'")
	
	// 测试时刻表相关端点
	testEndpoint(client, "每日时刻表 (v2)", "https://tdx.transportdata.tw/api/basic/v2/Rail/TRA/DailyTrainTimetable")
	testEndpoint(client, "每日时刻表 (v3)", "https://tdx.transportdata.tw/api/basic/v3/Rail/TRA/DailyTrainTimetable")
	
	// 测试其他可能的时刻表端点
	testEndpoint(client, "一般时刻表", "https://tdx.transportdata.tw/api/basic/v2/Rail/TRA/GeneralTimetable")
	testEndpoint(client, "列车信息", "https://tdx.transportdata.tw/api/basic/v2/Rail/TRA/GeneralTrainInfo")
	
	// 测试实时信息
	testEndpoint(client, "车站实时信息", "https://tdx.transportdata.tw/api/basic/v3/Rail/TRA/StationLiveBoard")
	
	// 测试具体站点的实时信息
	testEndpoint(client, "竹北站实时信息", "https://tdx.transportdata.tw/api/basic/v3/Rail/TRA/StationLiveBoard/1180")
}