package api

import (
	"fmt"
	"tg-bot-demo/handlers"
	"time"

	"github.com/gin-gonic/gin"
)

// SetupRouter 設定 API 路由
func SetupRouter(botHandlers *handlers.BotHandlers, debug bool) *gin.Engine {
	// 設定 Gin 模式
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// 中介軟體
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(CORSMiddleware())

	// 建立 API 處理器
	apiHandlers := NewAPIHandlers(botHandlers)

	// 健康檢查
	r.GET("/health", apiHandlers.Health)

	// API v1 路由群組
	v1 := r.Group("/api/v1")
	{
		// 發送訊息相關
		v1.POST("/send", apiHandlers.SendMessage)
		v1.POST("/broadcast", apiHandlers.Broadcast)
		
		// 查詢相關
		v1.GET("/chats", apiHandlers.GetChats)
	}

	return r
}

// CORSMiddleware CORS 中介軟體
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// LoggerMiddleware 自定義日誌中介軟體
func LoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC3339),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}
