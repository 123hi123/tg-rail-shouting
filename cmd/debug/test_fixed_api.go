package main

import (
	"fmt"
	"log"
	"tg-rail-shouting/internal/tdx"
)

func main() {
	fmt.Println("Testing fixed TDX API...")
	
	client := tdx.NewClient("", "", "https://tdx.transportdata.tw/api/basic/v3", "")
	
	// 测试1: 获取竹北站信息
	fmt.Println("1. Testing station info for Zhubei (1180)...")
	station, err := client.GetStationInfo("1180")
	if err != nil {
		log.Printf("Error getting station info: %v", err)
	} else {
		fmt.Printf("Station found: %s (%s)\n", station.StationName.ZhTw, station.StationID)
	}
	
	// 测试2: 获取竹北站实时列车信息 (北上 direction=1)
	fmt.Println("\n2. Testing live train data for Zhubei (direction=1, 北上)...")
	trains, err := client.GetTrainTimetable("1180", 1)
	if err != nil {
		log.Printf("Error getting timetable: %v", err)
	} else {
		fmt.Printf("Found %d northbound trains:\n", len(trains))
		for i, train := range trains {
			if i >= 3 {
				break
			}
			fmt.Printf("- %s次 (%s) 到达: %s, 终点: %s\n", 
				train.TrainNo, train.TrainType, train.ArrivalTime, train.EndStation)
		}
	}
	
	// 测试3: 获取竹北站实时列车信息 (南下 direction=0)
	fmt.Println("\n3. Testing live train data for Zhubei (direction=0, 南下)...")
	trains, err = client.GetTrainTimetable("1180", 0)
	if err != nil {
		log.Printf("Error getting timetable: %v", err)
	} else {
		fmt.Printf("Found %d southbound trains:\n", len(trains))
		for i, train := range trains {
			if i >= 3 {
				break
			}
			fmt.Printf("- %s次 (%s) 到达: %s, 终点: %s\n", 
				train.TrainNo, train.TrainType, train.ArrivalTime, train.EndStation)
		}
	}
	
	// 测试4: 如果实时数据没有，尝试获取一般时刻表
	fmt.Println("\n4. Testing general timetable (backup method)...")
	trains, err = client.GetGeneralTimetable("1180", 1)
	if err != nil {
		log.Printf("Error getting general timetable: %v", err)
	} else {
		fmt.Printf("Found %d trains in general timetable:\n", len(trains))
		for i, train := range trains {
			if i >= 2 {
				break
			}
			fmt.Printf("- %s次 (%s) 到达: %s, 途经 %d 站\n", 
				train.TrainNo, train.TrainType, train.ArrivalTime, len(train.Stations))
		}
	}
}