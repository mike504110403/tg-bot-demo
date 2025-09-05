package api

import (
	"net/http"
	"strconv"
	"tg-bot-demo/handlers"

	"github.com/gin-gonic/gin"
)

type APIHandlers struct {
	botHandlers *handlers.BotHandlers
}

// SendMessageRequest API 請求結構
type SendMessageRequest struct {
	ChatID  string `json:"chat_id" binding:"required"`  // 聊天ID（字串格式，支援數字和用戶名）
	Message string `json:"message" binding:"required"`  // 要發送的訊息內容
}

// SendMessageResponse API 回應結構
type SendMessageResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	ChatID  int64  `json:"chat_id,omitempty"`
}

// BroadcastRequest 廣播請求結構
type BroadcastRequest struct {
	Message string `json:"message" binding:"required"` // 要廣播的訊息內容
}

// BroadcastResponse 廣播回應結構
type BroadcastResponse struct {
	Success      bool `json:"success"`
	Message      string `json:"message"`
	SuccessCount int `json:"success_count"`
	FailCount    int `json:"fail_count"`
}

func NewAPIHandlers(botHandlers *handlers.BotHandlers) *APIHandlers {
	return &APIHandlers{
		botHandlers: botHandlers,
	}
}

// SendMessage 發送訊息到指定聊天
func (a *APIHandlers) SendMessage(c *gin.Context) {
	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, SendMessageResponse{
			Success: false,
			Message: "請求格式錯誤：" + err.Error(),
		})
		return
	}

	// 解析聊天ID
	chatID, err := strconv.ParseInt(req.ChatID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, SendMessageResponse{
			Success: false,
			Message: "無效的聊天ID格式，請使用數字格式",
		})
		return
	}

	// 發送訊息
	err = a.botHandlers.SendMessageToChat(chatID, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, SendMessageResponse{
			Success: false,
			Message: "發送訊息失敗：" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SendMessageResponse{
		Success: true,
		Message: "訊息發送成功",
		ChatID:  chatID,
	})
}

// Broadcast 廣播訊息到所有已知聊天
func (a *APIHandlers) Broadcast(c *gin.Context) {
	var req BroadcastRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, BroadcastResponse{
			Success: false,
			Message: "請求格式錯誤：" + err.Error(),
		})
		return
	}

	// 取得所有已知聊天ID
	chatIDs := a.botHandlers.GetKnownChatIDs()
	if len(chatIDs) == 0 {
		c.JSON(http.StatusBadRequest, BroadcastResponse{
			Success: false,
			Message: "目前沒有已知的聊天。請先讓機器人加入一些群組或與用戶互動。",
		})
		return
	}

	successCount := 0
	failCount := 0

	// 廣播訊息
	for _, chatID := range chatIDs {
		err := a.botHandlers.SendMessageToChat(chatID, req.Message)
		if err != nil {
			failCount++
		} else {
			successCount++
		}
	}

	c.JSON(http.StatusOK, BroadcastResponse{
		Success:      true,
		Message:      "廣播完成",
		SuccessCount: successCount,
		FailCount:    failCount,
	})
}

// GetChats 取得所有已知聊天ID
func (a *APIHandlers) GetChats(c *gin.Context) {
	chatIDs := a.botHandlers.GetKnownChatIDs()
	
	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"message":   "取得聊天列表成功",
		"chat_ids":  chatIDs,
		"count":     len(chatIDs),
	})
}

// Health 健康檢查端點
func (a *APIHandlers) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Telegram Bot API 運行正常",
	})
}
