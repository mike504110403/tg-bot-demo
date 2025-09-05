package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramBotToken string
	Debug            bool
	WebhookURL       string
	APIEnabled       bool
	APIPort          string
}

func Load() *Config {
	// 載入 .env 檔案（如果存在）
	if err := godotenv.Load(); err != nil {
		log.Printf("警告：無法載入 .env 檔案: %v", err)
	}

	// 從環境變數讀取配置
	config := &Config{
		TelegramBotToken: getEnv("TELEGRAM_BOT_TOKEN", ""),
		Debug:            getEnvAsBool("DEBUG", true),
		WebhookURL:       getEnv("WEBHOOK_URL", ""),
		APIEnabled:       getEnvAsBool("API_ENABLED", true),
		APIPort:          getEnv("API_PORT", "8080"),
	}

	// 驗證必要的配置
	if config.TelegramBotToken == "" {
		log.Fatal("錯誤：未設定 TELEGRAM_BOT_TOKEN 環境變數")
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
