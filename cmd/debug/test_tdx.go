package main

import (
	"fmt"
	"log"
	"tg-rail-shouting/internal/tdx"
)

func main() {
	fmt.Println("Testing TDX API without authentication...")
	
	client := tdx.NewClient("", "", "https://tdx.transportdata.tw/api/basic/v2", "")
	
	fmt.Println("Getting station info for Zhubei (1180)...")
	station, err := client.GetStationInfo("1180")
	if err != nil {
		log.Printf("Error getting station info: %v", err)
	} else {
		fmt.Printf("Station found: %s (%s)\n", station.StationName, station.StationID)
	}
	
	fmt.Println("\nGetting train timetable for Zhubei...")
	trains, err := client.GetTrainTimetable("1180", 1)
	if err != nil {
		log.Printf("Error getting timetable: %v", err)
	} else {
		fmt.Printf("Found %d trains:\n", len(trains))
		for i, train := range trains {
			if i >= 3 {
				break
			}
			fmt.Printf("- %s次 (%s) 到达: %s\n", train.TrainNo, train.TrainType, train.ArrivalTime)
		}
	}
}