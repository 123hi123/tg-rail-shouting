package tdx

import "time"

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type StationName struct {
	ZhTw string `json:"Zh_tw"`
	En   string `json:"En"`
}

type Station struct {
	StationUID  string      `json:"StationUID"`
	StationID   string      `json:"StationID"`
	StationName StationName `json:"StationName"`
	StationLat  float64     `json:"StationLat"`
	StationLon  float64     `json:"StationLon"`
}

type GeneralTimetableData struct {
	UpdateTime       string           `json:"UpdateTime"`
	VersionID        int              `json:"VersionID"`
	GeneralTimetable GeneralTimetable `json:"GeneralTimetable"`
}

type GeneralTimetable struct {
	GeneralTrainInfo GeneralTrainInfo `json:"GeneralTrainInfo"`
	StopTimes        []StopTime       `json:"StopTimes"`
	ServiceDay       ServiceDay       `json:"ServiceDay"`
}

type GeneralTrainInfo struct {
	TrainNo              string      `json:"TrainNo"`
	Direction            int         `json:"Direction"`
	StartingStationID    string      `json:"StartingStationID"`
	StartingStationName  StationName `json:"StartingStationName"`
	EndingStationID      string      `json:"EndingStationID"`
	EndingStationName    StationName `json:"EndingStationName"`
	TrainTypeID          string      `json:"TrainTypeID"`
	TrainTypeCode        string      `json:"TrainTypeCode"`
	TrainTypeName        StationName `json:"TrainTypeName"`
	TripLine             int         `json:"TripLine"`
	WheelchairFlag       int         `json:"WheelchairFlag"`
	PackageServiceFlag   int         `json:"PackageServiceFlag"`
	DiningFlag           int         `json:"DiningFlag"`
	BikeFlag             int         `json:"BikeFlag"`
	BreastFeedingFlag    int         `json:"BreastFeedingFlag"`
	DailyFlag            int         `json:"DailyFlag"`
	Note                 StationName `json:"Note"`
}

type StopTime struct {
	StopSequence  int         `json:"StopSequence"`
	StationID     string      `json:"StationID"`
	StationName   StationName `json:"StationName"`
	ArrivalTime   string      `json:"ArrivalTime"`
	DepartureTime string      `json:"DepartureTime"`
}

type ServiceDay struct {
	Monday    int `json:"Monday"`
	Tuesday   int `json:"Tuesday"`
	Wednesday int `json:"Wednesday"`
	Thursday  int `json:"Thursday"`
	Friday    int `json:"Friday"`
	Saturday  int `json:"Saturday"`
	Sunday    int `json:"Sunday"`
}

type StationLiveBoardResponse struct {
	UpdateTime         string              `json:"UpdateTime"`
	UpdateInterval     int                 `json:"UpdateInterval"`
	SrcUpdateTime      string              `json:"SrcUpdateTime"`
	SrcUpdateInterval  int                 `json:"SrcUpdateInterval"`
	AuthorityCode      string              `json:"AuthorityCode"`
	StationLiveBoards  []StationLiveBoard  `json:"StationLiveBoards"`
}

type StationLiveBoard struct {
	StationID              string      `json:"StationID"`
	StationName            StationName `json:"StationName"`
	TrainNo                string      `json:"TrainNo"`
	Direction              int         `json:"Direction"`
	TrainTypeID            string      `json:"TrainTypeID"`
	TrainTypeCode          string      `json:"TrainTypeCode"`
	TrainTypeName          StationName `json:"TrainTypeName"`
	EndingStationID        string      `json:"EndingStationID"`
	EndingStationName      StationName `json:"EndingStationName"`
	TripLine               int         `json:"TripLine"`
	Platform               string      `json:"Platform"`
	ScheduleArrivalTime    string      `json:"ScheduleArrivalTime"`
	ScheduleDepartureTime  string      `json:"ScheduleDepartureTime"`
	DelayTime              int         `json:"DelayTime"`
	RunningStatus          int         `json:"RunningStatus"`
	UpdateTime             time.Time   `json:"UpdateTime"`
}

// 简化的数据结构，用于应用逻辑
type TrainInfo struct {
	TrainNo       string
	TrainType     string
	ArrivalTime   string
	DepartureTime string
	StopSequence  int
	Stations      []StationInfo
	Direction     int
	EndStation    string
}

type StationInfo struct {
	StationName   string
	ArrivalTime   string
	DepartureTime string
	StopSequence  int
}