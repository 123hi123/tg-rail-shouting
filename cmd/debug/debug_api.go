package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

func main() {
	client := resty.New()
	client.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	
	fmt.Println("Testing different TDX API endpoints...")
	
	// Test 1: Station info
	fmt.Println("\n1. Testing station info endpoint...")
	resp, err := client.R().
		SetQueryParam("$format", "JSON").
		SetQueryParam("$filter", "StationID eq '1180'").
		Get("https://tdx.transportdata.tw/api/basic/v2/Rail/TRA/Station")
	
	fmt.Printf("Status: %d\n", resp.StatusCode())
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else if resp.StatusCode() != 200 {
		fmt.Printf("Response body: %s\n", resp.String())
	} else {
		fmt.Printf("Success! Response length: %d\n", len(resp.Body()))
	}
	
	// Test 2: Try different timetable endpoints
	endpoints := []string{
		"https://tdx.transportdata.tw/api/basic/v2/Rail/TRA/DailyTrainTimetable/Station/1180",
		"https://tdx.transportdata.tw/api/basic/v2/Rail/TRA/GeneralTimetable/Station/1180", 
		"https://tdx.transportdata.tw/api/basic/v3/Rail/TRA/StationLiveBoard/1180",
		"https://tdx.transportdata.tw/api/basic/v2/Rail/TRA/GeneralTrainInfo",
	}
	
	for i, endpoint := range endpoints {
		fmt.Printf("\n%d. Testing endpoint: %s\n", i+2, endpoint)
		resp, err := client.R().
			SetQueryParam("$format", "JSON").
			Get(endpoint)
		
		fmt.Printf("Status: %d\n", resp.StatusCode())
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else if resp.StatusCode() != 200 {
			fmt.Printf("Response body: %s\n", resp.String())
		} else {
			fmt.Printf("Success! Response length: %d\n", len(resp.Body()))
		}
	}
}