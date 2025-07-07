package tdx

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

type Client struct {
	client       *resty.Client
	clientID     string
	clientSecret string
	baseURL      string
	authURL      string
	accessToken  string
	tokenExpiry  time.Time
}

func NewClient(clientID, clientSecret, baseURL, authURL string) *Client {
	client := resty.New()
	client.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	
	return &Client{
		client:       client,
		clientID:     clientID,
		clientSecret: clientSecret,
		baseURL:      baseURL,
		authURL:      authURL,
	}
}

func (c *Client) authenticate() error {
	// 如果没有提供认证信息，使用免费API
	if c.clientID == "" || c.clientSecret == "" {
		logrus.Info("Using TDX API without authentication (free tier)")
		return nil
	}
	
	if time.Now().Before(c.tokenExpiry) {
		return nil
	}

	logrus.Info("Authenticating with TDX API...")
	
	resp, err := c.client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"grant_type":    "client_credentials",
			"client_id":     c.clientID,
			"client_secret": c.clientSecret,
		}).
		Post(c.authURL)

	if err != nil {
		return fmt.Errorf("authentication request failed: %w", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("authentication failed with status: %d, body: %s", resp.StatusCode(), resp.String())
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(resp.Body(), &tokenResp); err != nil {
		return fmt.Errorf("failed to parse token response: %w", err)
	}

	c.accessToken = tokenResp.AccessToken
	c.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn-60) * time.Second)
	
	logrus.Info("TDX API authentication successful")
	return nil
}

func (c *Client) GetStationInfo(stationID string) (*Station, error) {
	if err := c.authenticate(); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/Rail/TRA/Station", c.baseURL)
	filter := fmt.Sprintf("StationID eq '%s'", stationID)

	req := c.client.R().
		SetQueryParam("$filter", filter).
		SetQueryParam("$format", "JSON")
	
	if c.accessToken != "" {
		req.SetHeader("Authorization", "Bearer "+c.accessToken)
	}
	
	resp, err := req.Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to get station info: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("API request failed with status: %d, body: %s", resp.StatusCode(), resp.String())
	}

	var stations []Station
	if err := json.Unmarshal(resp.Body(), &stations); err != nil {
		return nil, fmt.Errorf("failed to parse station response: %w", err)
	}

	if len(stations) == 0 {
		return nil, fmt.Errorf("station not found: %s", stationID)
	}

	return &stations[0], nil
}

// GetTrainTimetable 获取车站的实时列车信息
func (c *Client) GetTrainTimetable(stationID string, direction int) ([]TrainInfo, error) {
	if err := c.authenticate(); err != nil {
		return nil, err
	}

	// 使用StationLiveBoard获取实时信息
	url := fmt.Sprintf("%s/Rail/TRA/StationLiveBoard", c.baseURL)
	filter := fmt.Sprintf("StationID eq '%s'", stationID)
	
	req := c.client.R().
		SetQueryParam("$format", "JSON").
		SetQueryParam("$filter", filter)
	
	if c.accessToken != "" {
		req.SetHeader("Authorization", "Bearer "+c.accessToken)
	}
	
	resp, err := req.Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to get station live board: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("API request failed with status: %d, body: %s", resp.StatusCode(), resp.String())
	}

	var liveBoard StationLiveBoardResponse
	if err := json.Unmarshal(resp.Body(), &liveBoard); err != nil {
		return nil, fmt.Errorf("failed to parse live board response: %w", err)
	}

	var trains []TrainInfo
	now := time.Now()
	currentTime := now.Format("15:04:05")

	for _, board := range liveBoard.StationLiveBoards {
		if board.StationID == stationID && board.Direction == direction {
			arrivalTime := board.ScheduleArrivalTime
			if arrivalTime == "" {
				arrivalTime = board.ScheduleDepartureTime
			}

			// 只显示当前时间之后的列车
			if arrivalTime >= currentTime {
				trainInfo := TrainInfo{
					TrainNo:       board.TrainNo,
					TrainType:     board.TrainTypeName.ZhTw,
					ArrivalTime:   arrivalTime,
					DepartureTime: board.ScheduleDepartureTime,
					Direction:     board.Direction,
					EndStation:    board.EndingStationName.ZhTw,
				}

				trains = append(trains, trainInfo)
			}
		}
	}

	// 按到达时间排序
	sort.Slice(trains, func(i, j int) bool {
		return trains[i].ArrivalTime < trains[j].ArrivalTime
	})

	return trains, nil
}

// GetGeneralTimetable 获取完整的时刻表数据（备用方法）
func (c *Client) GetGeneralTimetable(stationID string, direction int) ([]TrainInfo, error) {
	if err := c.authenticate(); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/Rail/TRA/GeneralTimetable", c.baseURL)
	
	req := c.client.R().
		SetQueryParam("$format", "JSON").
		SetQueryParam("$top", "100")
	
	if c.accessToken != "" {
		req.SetHeader("Authorization", "Bearer "+c.accessToken)
	}
	
	resp, err := req.Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to get general timetable: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("API request failed with status: %d, body: %s", resp.StatusCode(), resp.String())
	}

	var timetables []GeneralTimetableData
	if err := json.Unmarshal(resp.Body(), &timetables); err != nil {
		return nil, fmt.Errorf("failed to parse timetable response: %w", err)
	}

	var trains []TrainInfo
	now := time.Now()
	currentTime := now.Format("15:04:05")

	for _, tt := range timetables {
		if tt.GeneralTimetable.GeneralTrainInfo.Direction != direction {
			continue
		}

		for _, st := range tt.GeneralTimetable.StopTimes {
			if st.StationID == stationID {
				arrivalTime := st.ArrivalTime
				if arrivalTime == "" {
					arrivalTime = st.DepartureTime
				}

				if arrivalTime >= currentTime {
					trainInfo := TrainInfo{
						TrainNo:       tt.GeneralTimetable.GeneralTrainInfo.TrainNo,
						TrainType:     tt.GeneralTimetable.GeneralTrainInfo.TrainTypeName.ZhTw,
						ArrivalTime:   arrivalTime,
						DepartureTime: st.DepartureTime,
						StopSequence:  st.StopSequence,
						Direction:     tt.GeneralTimetable.GeneralTrainInfo.Direction,
						EndStation:    tt.GeneralTimetable.GeneralTrainInfo.EndingStationName.ZhTw,
					}

					stations := c.extractStationInfo(tt.GeneralTimetable.StopTimes, st.StopSequence)
					trainInfo.Stations = stations

					trains = append(trains, trainInfo)
				}
				break
			}
		}
	}

	sort.Slice(trains, func(i, j int) bool {
		return trains[i].ArrivalTime < trains[j].ArrivalTime
	})

	return trains, nil
}

func (c *Client) extractStationInfo(stopTimes []StopTime, currentSequence int) []StationInfo {
	var stations []StationInfo
	
	for _, st := range stopTimes {
		if st.StopSequence >= currentSequence {
			stationInfo := StationInfo{
				StationName:   st.StationName.ZhTw,
				ArrivalTime:   st.ArrivalTime,
				DepartureTime: st.DepartureTime,
				StopSequence:  st.StopSequence,
			}
			stations = append(stations, stationInfo)
		}
	}
	
	return stations
}

func (c *Client) GetTrainRoute(trainNo string) ([]StationInfo, error) {
	if err := c.authenticate(); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/Rail/TRA/GeneralTimetable", c.baseURL)
	filter := fmt.Sprintf("GeneralTimetable/GeneralTrainInfo/TrainNo eq '%s'", trainNo)
	
	req := c.client.R().
		SetQueryParam("$format", "JSON").
		SetQueryParam("$filter", filter)
	
	if c.accessToken != "" {
		req.SetHeader("Authorization", "Bearer "+c.accessToken)
	}
	
	resp, err := req.Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to get train route: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("API request failed with status: %d, body: %s", resp.StatusCode(), resp.String())
	}

	var timetables []GeneralTimetableData
	if err := json.Unmarshal(resp.Body(), &timetables); err != nil {
		return nil, fmt.Errorf("failed to parse route response: %w", err)
	}

	if len(timetables) == 0 {
		return nil, fmt.Errorf("train not found: %s", trainNo)
	}

	var stations []StationInfo
	for _, st := range timetables[0].GeneralTimetable.StopTimes {
		stationInfo := StationInfo{
			StationName:   st.StationName.ZhTw,
			ArrivalTime:   st.ArrivalTime,
			DepartureTime: st.DepartureTime,
			StopSequence:  st.StopSequence,
		}
		stations = append(stations, stationInfo)
	}

	return stations, nil
}

func (c *Client) FindRouteToFugang(trainNo string, fromStationID string) ([]StationInfo, bool, error) {
	route, err := c.GetTrainRoute(trainNo)
	if err != nil {
		return nil, false, err
	}

	var fromIndex = -1
	var fugangIndex = -1
	
	for i, station := range route {
		if strings.Contains(station.StationName, "竹北") || 
		   strings.Contains(station.StationName, fromStationID) {
			fromIndex = i
		}
		if strings.Contains(station.StationName, "富岡") {
			fugangIndex = i
		}
	}

	if fromIndex == -1 {
		return route, false, nil
	}

	reachFugang := fugangIndex != -1 && fugangIndex > fromIndex
	
	if reachFugang {
		return route[fromIndex:fugangIndex+1], true, nil
	}

	return route[fromIndex:], false, nil
}