package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"tg-rail-shouting/internal/config"
	"tg-rail-shouting/internal/monitor"
	"tg-rail-shouting/internal/tdx"
	"tg-rail-shouting/internal/telegram"
)

func main() {
	setupLogger()
	
	logrus.Info("Starting TG Rail Shouting service...")
	
	cfg, err := config.Load()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load configuration")
	}
	
	tdxClient := tdx.NewClient(
		cfg.TDX.ClientID,
		cfg.TDX.ClientSecret,
		cfg.TDX.BaseURL,
		cfg.TDX.AuthURL,
	)
	
	tgBot := telegram.NewBot(cfg.Telegram.BotToken, cfg.Telegram.ChatID)
	
	scheduler := monitor.NewScheduler(cfg, tdxClient, tgBot)
	
	if err := scheduler.SendTestMessage(); err != nil {
		logrus.WithError(err).Warn("Failed to send test message")
	}
	
	if err := scheduler.Start(); err != nil {
		logrus.WithError(err).Fatal("Failed to start scheduler")
	}
	
	logrus.Info("Service started successfully")
	
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	
	<-stop
	logrus.Info("Received shutdown signal")
	
	scheduler.Stop()
	logrus.Info("Service stopped")
}

func setupLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logrus.SetLevel(logrus.InfoLevel)
}