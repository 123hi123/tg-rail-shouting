package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	TDX TDXConfig
	Telegram TelegramConfig  
	Monitor MonitorConfig
	Station StationConfig
}

type TDXConfig struct {
	ClientID     string
	ClientSecret string
	BaseURL      string
	AuthURL      string
}

type TelegramConfig struct {
	BotToken string
	ChatID   string
}

type MonitorConfig struct {
	StartHour        int
	EndHour          int
	IntervalMinutes  int
}

type StationConfig struct {
	ZhubeiStationID  string
	TargetDirection  int
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No .env file found, using environment variables")
	}

	config := &Config{
		TDX: TDXConfig{
			ClientID:     os.Getenv("TDX_CLIENT_ID"),
			ClientSecret: os.Getenv("TDX_CLIENT_SECRET"),
			BaseURL:      "https://tdx.transportdata.tw/api/basic/v3",
			AuthURL:      "https://tdx.transportdata.tw/auth/realms/TDXConnect/protocol/openid-connect/token",
		},
		Telegram: TelegramConfig{
			BotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
			ChatID:   os.Getenv("TELEGRAM_CHAT_ID"),
		},
		Monitor: MonitorConfig{
			StartHour:       getIntEnv("MONITOR_START_HOUR", 18),
			EndHour:         getIntEnv("MONITOR_END_HOUR", 23),
			IntervalMinutes: getIntEnv("MONITOR_INTERVAL_MINUTES", 30),
		},
		Station: StationConfig{
			ZhubeiStationID: os.Getenv("ZHUBEI_STATION_ID"),
			TargetDirection: getIntEnv("TARGET_DIRECTION", 1),
		},
	}

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func validateConfig(config *Config) error {
	// TDX 认证信息可选（使用免费API）
	if config.TDX.ClientID == "" || config.TDX.ClientSecret == "" {
		logrus.Warn("TDX API credentials not provided, using free tier (50 requests/day limit)")
	}
	
	if config.Telegram.BotToken == "" {
		logrus.Fatal("TELEGRAM_BOT_TOKEN is required")
	}
	if config.Telegram.ChatID == "" {
		logrus.Fatal("TELEGRAM_CHAT_ID is required")
	}
	if config.Station.ZhubeiStationID == "" {
		logrus.Fatal("ZHUBEI_STATION_ID is required")
	}
	return nil
}