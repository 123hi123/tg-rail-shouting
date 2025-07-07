package monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"tg-rail-shouting/internal/config"
	"tg-rail-shouting/internal/tdx"
	"tg-rail-shouting/internal/telegram"
)

type Scheduler struct {
	cron      *cron.Cron
	config    *config.Config
	tdxClient *tdx.Client
	tgBot     *telegram.Bot
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewScheduler(cfg *config.Config, tdxClient *tdx.Client, tgBot *telegram.Bot) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &Scheduler{
		cron:      cron.New(cron.WithLocation(time.Local)),
		config:    cfg,
		tdxClient: tdxClient,
		tgBot:     tgBot,
		ctx:       ctx,
		cancel:    cancel,
	}
}

func (s *Scheduler) Start() error {
	cronExpr := fmt.Sprintf("*/%d * * * *", s.config.Monitor.IntervalMinutes)
	
	_, err := s.cron.AddFunc(cronExpr, func() {
		if s.shouldMonitor() {
			s.checkTrains()
		}
	})
	
	if err != nil {
		return fmt.Errorf("failed to add cron job: %w", err)
	}
	
	s.cron.Start()
	logrus.Info("Scheduler started")
	
	go s.runInitialCheck()
	
	return nil
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
	s.cancel()
	logrus.Info("Scheduler stopped")
}

func (s *Scheduler) shouldMonitor() bool {
	now := time.Now()
	hour := now.Hour()
	return hour >= s.config.Monitor.StartHour && hour <= s.config.Monitor.EndHour
}

func (s *Scheduler) runInitialCheck() {
	time.Sleep(3 * time.Second)
	
	logrus.Info("Running initial train check to verify service...")
	s.checkTrainsForce(true)
}

func (s *Scheduler) checkTrains() {
	if s.shouldMonitor() {
		s.checkTrainsForce(false)
	}
}

func (s *Scheduler) checkTrainsForce(isInitial bool) {
	select {
	case <-s.ctx.Done():
		return
	default:
	}
	
	if isInitial {
		logrus.Info("Initial API test - checking trains for Zhubei Station...")
	} else {
		logrus.Info("Scheduled check - checking trains for Zhubei Station...")
	}
	
	trains, err := s.tdxClient.GetTrainTimetable(s.config.Station.ZhubeiStationID, s.config.Station.TargetDirection)
	if err != nil {
		logrus.WithError(err).Error("Failed to get train timetable")
		if isInitial {
			s.sendInitialErrorMessage(err)
		} else {
			s.sendErrorMessage(err)
		}
		return
	}
	
	if len(trains) == 0 {
		logrus.Info("No trains found for current time")
		if isInitial {
			s.sendNoTrainsMessage()
		}
		return
	}
	
	s.processTrains(trains, isInitial)
}

func (s *Scheduler) processTrains(trains []tdx.TrainInfo, isInitial bool) {
	now := time.Now()
	currentTime := now.Format("15:04:05")
	
	var upcomingTrains []tdx.TrainInfo
	for _, train := range trains {
		if train.ArrivalTime >= currentTime {
			upcomingTrains = append(upcomingTrains, train)
		}
	}
	
	if len(upcomingTrains) == 0 {
		logrus.Info("No upcoming trains found")
		if isInitial {
			s.sendNoTrainsMessage()
		}
		return
	}
	
	logrus.WithField("count", len(upcomingTrains)).Info("Found upcoming trains")
	
	var stationName string
	if isInitial {
		stationName = "竹北 (服务测试)"
	} else {
		stationName = "竹北"
	}
	
	if err := s.tgBot.SendTrainInfo(upcomingTrains, stationName); err != nil {
		logrus.WithError(err).Error("Failed to send train info")
		return
	}
	
	if !isInitial {
		s.sendDetailedInfo(upcomingTrains)
	}
}

func (s *Scheduler) sendDetailedInfo(trains []tdx.TrainInfo) {
	if len(trains) == 0 {
		return
	}
	
	var trainsToFugang []tdx.TrainInfo
	
	for _, train := range trains {
		if len(trainsToFugang) >= 3 {
			break
		}
		
		route, reachFugang, err := s.tdxClient.FindRouteToFugang(train.TrainNo, s.config.Station.ZhubeiStationID)
		if err != nil {
			logrus.WithError(err).WithField("train", train.TrainNo).Warn("Failed to get route to Fugang")
			continue
		}
		
		if reachFugang {
			trainWithRoute := train
			trainWithRoute.Stations = make([]tdx.StationInfo, len(route))
			for i, station := range route {
				trainWithRoute.Stations[i] = tdx.StationInfo{
					StationName:   station.StationName,
					ArrivalTime:   station.ArrivalTime,
					DepartureTime: station.DepartureTime,
					StopSequence:  station.StopSequence,
				}
			}
			trainsToFugang = append(trainsToFugang, trainWithRoute)
		}
	}
	
	if len(trainsToFugang) > 0 {
		if err := s.tgBot.SendDetailedTrainInfo(trainsToFugang, "竹北", "富岡"); err != nil {
			logrus.WithError(err).Error("Failed to send detailed train info")
		}
	}
}

func (s *Scheduler) sendErrorMessage(err error) {
	message := fmt.Sprintf("❌ 获取列车信息失败\n\n错误: %v\n时间: %s", err, time.Now().Format("2006-01-02 15:04:05"))
	
	if sendErr := s.tgBot.SendMessage(message); sendErr != nil {
		logrus.WithError(sendErr).Error("Failed to send error message")
	}
}

func (s *Scheduler) sendInitialErrorMessage(err error) {
	message := fmt.Sprintf("⚠️ 台铁监控服务已启动，但API测试失败\n\n监控时间: %d:00 - %d:00\n检查间隔: %d分钟\n监控站点: 竹北站\n\n❌ API测试错误: %v\n时间: %s\n\n服务将继续运行，稍后会重试...", 
		s.config.Monitor.StartHour, 
		s.config.Monitor.EndHour,
		s.config.Monitor.IntervalMinutes,
		err, 
		time.Now().Format("2006-01-02 15:04:05"))
	
	if sendErr := s.tgBot.SendMessage(message); sendErr != nil {
		logrus.WithError(sendErr).Error("Failed to send initial error message")
	}
}

func (s *Scheduler) sendNoTrainsMessage() {
	now := time.Now()
	message := fmt.Sprintf("✅ 台铁监控服务已启动并完成API测试\n\n监控时间: %d:00 - %d:00\n检查间隔: %d分钟\n监控站点: 竹北站\n\n🚄 API测试结果: 当前时间(%s)没有列车信息\n这很正常，服务将在监控时间内定期检查\n\n服务运行正常 ✅", 
		s.config.Monitor.StartHour, 
		s.config.Monitor.EndHour,
		s.config.Monitor.IntervalMinutes,
		now.Format("15:04"))
	
	if sendErr := s.tgBot.SendMessage(message); sendErr != nil {
		logrus.WithError(sendErr).Error("Failed to send no trains message")
	}
}

func (s *Scheduler) SendTestMessage() error {
	message := fmt.Sprintf("✅ 台铁监控服务已启动\n\n监控时间: %d:00 - %d:00\n检查间隔: %d分钟\n监控站点: 竹北站\n\n正在进行API连接测试...", 
		s.config.Monitor.StartHour, 
		s.config.Monitor.EndHour,
		s.config.Monitor.IntervalMinutes)
	
	return s.tgBot.SendMessage(message)
}