package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"tg-bot-demo/api"
	"tg-bot-demo/config"
	"tg-bot-demo/handlers"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// 載入配置
	cfg := config.Load()

	// 建立機器人實例
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		log.Fatalf("建立機器人時發生錯誤: %v", err)
	}

	// 設定除錯模式
	bot.Debug = cfg.Debug

	log.Printf("已授權帳號 %s", bot.Self.UserName)

	// 建立處理程序
	botHandlers := handlers.NewBotHandlers(bot)

	// 使用 WaitGroup 來管理 goroutines
	var wg sync.WaitGroup

	// 建立 context 用於優雅關閉
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 啟動 Telegram Bot
	wg.Add(1)
	go func() {
		defer wg.Done()
		runBot(ctx, bot, botHandlers)
	}()

	// 啟動 HTTP API（如果啟用）
	if cfg.APIEnabled {
		wg.Add(1)
		go func() {
			defer wg.Done()
			runAPI(ctx, cfg, botHandlers)
		}()
	}

	// 等待中斷信號
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("收到關閉信號，正在優雅關閉...")

	cancel()  // 取消 context
	wg.Wait() // 等待所有 goroutines 結束

	log.Println("應用程式已關閉")
}

// runBot 運行 Telegram Bot
func runBot(ctx context.Context, bot *tgbotapi.BotAPI, botHandlers *handlers.BotHandlers) {
	log.Println("Telegram 機器人開始運行...")

	// 建立更新配置
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// 取得更新通道
	updates := bot.GetUpdatesChan(u)

	// 處理更新
	for {
		select {
		case <-ctx.Done():
			log.Println("Telegram 機器人正在關閉...")
			bot.StopReceivingUpdates()
			return
		case update := <-updates:
			// 在 goroutine 中處理每個更新，避免阻塞
			go botHandlers.HandleUpdate(update)
		}
	}
}

// runAPI 運行 HTTP API 伺服器
func runAPI(ctx context.Context, cfg *config.Config, botHandlers *handlers.BotHandlers) {
	// 設定路由
	router := api.SetupRouter(botHandlers, cfg.Debug)

	// 建立 HTTP 伺服器
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.APIPort),
		Handler: router,
	}

	// 在 goroutine 中啟動伺服器
	go func() {
		log.Printf("HTTP API 伺服器啟動於 http://localhost:%s", cfg.APIPort)
		log.Printf("API 文件：")
		log.Printf("  POST /api/v1/send - 發送訊息")
		log.Printf("  POST /api/v1/broadcast - 廣播訊息")
		log.Printf("  GET  /api/v1/chats - 取得聊天列表")
		log.Printf("  GET  /health - 健康檢查")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP 伺服器錯誤: %v", err)
		}
	}()

	// 等待 context 取消
	<-ctx.Done()

	log.Println("HTTP API 伺服器正在關閉...")

	// 優雅關閉伺服器
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("伺服器關閉錯誤: %v", err)
	} else {
		log.Println("HTTP API 伺服器已關閉")
	}
}
