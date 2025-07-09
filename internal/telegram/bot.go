package telegram

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"tg-rail-shouting/internal/tdx"
)

type Bot struct {
	client *resty.Client
	token  string
	chatID string
}

func NewBot(token, chatID string) *Bot {
	return &Bot{
		client: resty.New(),
		token:  token,
		chatID: chatID,
	}
}

func (b *Bot) SendMessage(text string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", b.token)
	
	resp, err := b.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"chat_id":    b.chatID,
			"text":       text,
			"parse_mode": "HTML",
		}).
		Post(url)

	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("telegram API error: %d, body: %s", resp.StatusCode(), resp.String())
	}

	logrus.Info("Message sent successfully to Telegram")
	return nil
}

func (b *Bot) SendTrainInfo(trains []tdx.TrainInfo, stationName string) error {
	if len(trains) == 0 {
		message := fmt.Sprintf("🚄 %s站 列车信息\n\n暂无列车信息", stationName)
		return b.SendMessage(message)
	}

	var message strings.Builder
	message.WriteString(fmt.Sprintf("🚄 %s站 列车信息\n\n", stationName))

	for i, train := range trains {
		if i >= 5 {
			break
		}

		message.WriteString(fmt.Sprintf("🚂 %s次 (%s)\n", train.TrainNo, train.TrainType))
		message.WriteString(fmt.Sprintf("⏰ 到达: %s", train.ArrivalTime))
		if train.DepartureTime != "" && train.DepartureTime != train.ArrivalTime {
			message.WriteString(fmt.Sprintf(" / 出发: %s", train.DepartureTime))
		}
		message.WriteString("\n\n")

		if len(train.Stations) > 0 {
			message.WriteString("完整路線: ")
			stationNames := make([]string, 0, len(train.Stations))
			for _, station := range train.Stations {
				stationName := station.StationName
				// 如果是富岡站，加粗顯示
				if strings.Contains(stationName, "富岡") {
					stationName = fmt.Sprintf("<b>%s</b>", stationName)
				}
				stationNames = append(stationNames, stationName)
			}
			message.WriteString(strings.Join(stationNames, " → "))
			message.WriteString("\n")
		}
		
		message.WriteString("\n")
	}

	return b.SendMessage(message.String())
}

func (b *Bot) SendDetailedTrainInfo(trains []tdx.TrainInfo, stationName string, targetStation string) error {
	if len(trains) == 0 {
		message := fmt.Sprintf("🚄 <b>%s站 → %s 列车信息</b>\n\n暂无列车信息", stationName, targetStation)
		return b.SendMessage(message)
	}

	var message strings.Builder
	message.WriteString(fmt.Sprintf("🚄 <b>%s站 → %s 列车信息</b>\n", stationName, targetStation))
	message.WriteString(fmt.Sprintf("📅 更新时间: %s\n\n", getCurrentTime()))

	for i, train := range trains {
		if i >= 5 {
			break
		}

		message.WriteString(fmt.Sprintf("🚂 <b>%s次 (%s)</b>\n", train.TrainNo, train.TrainType))
		message.WriteString(fmt.Sprintf("⏰ 到达: %s", train.ArrivalTime))
		if train.DepartureTime != "" && train.DepartureTime != train.ArrivalTime {
			message.WriteString(fmt.Sprintf(" / 出发: %s", train.DepartureTime))
		}
		message.WriteString("\n")

		if len(train.Stations) > 0 {
			message.WriteString("🛤️ 完整路线:\n")
			for j, station := range train.Stations {
				if j > 15 {
					message.WriteString("    ...\n")
					break
				}
				timeStr := station.ArrivalTime
				if timeStr == "" {
					timeStr = station.DepartureTime
				}
				
				emoji := "  "
				if strings.Contains(station.StationName, "竹北") {
					emoji = "🔵"
				} else if strings.Contains(station.StationName, "富岡") {
					emoji = "🔴"
				}
				
				message.WriteString(fmt.Sprintf("    %s %s (%s)\n", emoji, station.StationName, timeStr))
			}
		}
		
		message.WriteString("\n")
	}

	return b.SendMessage(message.String())
}

func (b *Bot) SendStartupMessage() error {
	version := b.getVersion()
	message := fmt.Sprintf("🚀 <b>台灣鐵路監控服務啟動成功</b>\n\n"+
		"✅ Telegram Bot 連線正常\n"+
		"✅ 服務配置載入完成\n"+
		"⏰ 監控時間: 18:00-23:00\n"+
		"🔄 檢查間隔: 每30分鐘\n\n"+
		"📋 版本: v%s", version)
	
	return b.SendMessage(message)
}

func (b *Bot) getVersion() string {
	content, err := os.ReadFile("version.txt")
	if err != nil {
		logrus.WithError(err).Warn("Failed to read version file")
		return "unknown"
	}
	return strings.TrimSpace(string(content))
}

func getCurrentTime() string {
	return fmt.Sprintf("%s", "现在")
}