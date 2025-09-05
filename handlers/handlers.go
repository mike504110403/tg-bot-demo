package handlers

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotHandlers struct {
	bot       *tgbotapi.BotAPI
	chatIDs   map[int64]bool // 儲存已知的聊天ID
	chatMutex sync.RWMutex   // 保護 chatIDs 的讀寫
}

func NewBotHandlers(bot *tgbotapi.BotAPI) *BotHandlers {
	return &BotHandlers{
		bot:     bot,
		chatIDs: make(map[int64]bool),
	}
}

// HandleUpdate 處理所有更新
func (h *BotHandlers) HandleUpdate(update tgbotapi.Update) {
	if update.Message != nil {
		// 記錄聊天ID
		h.addChatID(update.Message.Chat.ID)
		h.handleMessage(update.Message)
	} else if update.CallbackQuery != nil {
		h.handleCallbackQuery(update.CallbackQuery)
	}
}

// addChatID 添加聊天ID到已知列表
func (h *BotHandlers) addChatID(chatID int64) {
	h.chatMutex.Lock()
	defer h.chatMutex.Unlock()
	h.chatIDs[chatID] = true
}

// handleMessage 處理文字訊息
func (h *BotHandlers) handleMessage(message *tgbotapi.Message) {
	if message.IsCommand() {
		h.handleCommand(message)
		return
	}

	// 處理一般文字訊息
	h.handleTextMessage(message)
}

// handleCommand 處理指令
func (h *BotHandlers) handleCommand(message *tgbotapi.Message) {
	command := message.Command()
	args := message.CommandArguments()

	log.Printf("收到指令: /%s，參數: %s，來自用戶: %s", command, args, message.From.UserName)

	switch command {
	case "start":
		h.handleStartCommand(message)
	case "help":
		h.handleHelpCommand(message)
	case "info":
		h.handleInfoCommand(message)
	case "echo":
		h.handleEchoCommand(message, args)
	case "broadcast":
		h.handleBroadcastCommand(message, args)
	case "sendto":
		h.handleSendToCommand(message, args)
	default:
		h.sendReply(message, "抱歉，我不認識這個指令。輸入 /help 查看可用指令。")
	}
}

// handleStartCommand 處理 /start 指令
func (h *BotHandlers) handleStartCommand(message *tgbotapi.Message) {
	welcomeText := fmt.Sprintf(`👋 歡迎使用此機器人，%s！

我是一個用 Go 語言開發的 Telegram 機器人。

可用指令：
/start - 顯示歡迎訊息
/help - 顯示幫助訊息
/info - 顯示您的資訊
/echo <文字> - 重複您輸入的文字
/sendto <聊天ID> <訊息> - 發送訊息到指定聊天
/broadcast <訊息> - 廣播訊息到所有已知聊天

您也可以直接發送訊息給我，我會回應您！`, message.From.FirstName)

	h.sendReply(message, welcomeText)
}

// handleHelpCommand 處理 /help 指令
func (h *BotHandlers) handleHelpCommand(message *tgbotapi.Message) {
	helpText := `📚 可用指令說明：

/start - 顯示歡迎訊息
/help - 顯示此幫助訊息
/info - 顯示您的用戶資訊
/echo <文字> - 重複您輸入的文字
/sendto <聊天ID> <訊息> - 發送訊息到指定聊天
/broadcast <訊息> - 廣播訊息到所有已知聊天

💡 您也可以直接發送任何文字訊息，我會回應您！`

	h.sendReply(message, helpText)
}

// handleInfoCommand 處理 /info 指令
func (h *BotHandlers) handleInfoCommand(message *tgbotapi.Message) {
	user := message.From
	chat := message.Chat

	infoText := fmt.Sprintf(`ℹ️ 您的資訊：

👤 用戶名稱: %s %s
🆔 用戶 ID: %d
📝 用戶名: @%s
💬 聊天類型: %s
🔢 聊天 ID: %d`,
		user.FirstName,
		user.LastName,
		user.ID,
		user.UserName,
		chat.Type,
		chat.ID)

	h.sendReply(message, infoText)
}

// handleEchoCommand 處理 /echo 指令
func (h *BotHandlers) handleEchoCommand(message *tgbotapi.Message, args string) {
	if args == "" {
		h.sendReply(message, "請在 /echo 後面輸入要重複的文字。\n例如：/echo Hello World")
		return
	}

	echoText := fmt.Sprintf("🔄 您說：%s", args)
	h.sendReply(message, echoText)
}

// handleTextMessage 處理一般文字訊息
func (h *BotHandlers) handleTextMessage(message *tgbotapi.Message) {
	text := strings.TrimSpace(message.Text)
	log.Printf("收到文字訊息: %s，來自用戶: %s", text, message.From.UserName)

	// 簡單的回應邏輯
	var reply string
	switch strings.ToLower(text) {
	case "你好", "hello", "hi":
		reply = fmt.Sprintf("你好 %s！很高興見到你 😊", message.From.FirstName)
	case "謝謝", "thank you", "thanks":
		reply = "不客氣！很高興能幫助您 😊"
	case "再見", "bye", "goodbye":
		reply = "再見！期待下次見面 👋"
	default:
		reply = fmt.Sprintf("你發送了：「%s」\n\n我是一個簡單的機器人，還在學習中！輸入 /help 查看我能做什麼。", text)
	}

	h.sendReply(message, reply)
}

// handleCallbackQuery 處理回調查詢
func (h *BotHandlers) handleCallbackQuery(callbackQuery *tgbotapi.CallbackQuery) {
	// 回應回調查詢
	callback := tgbotapi.NewCallback(callbackQuery.ID, "已處理！")
	if _, err := h.bot.Request(callback); err != nil {
		log.Printf("回應回調查詢時發生錯誤: %v", err)
	}

	// 可以在這裡處理不同的回調數據
	data := callbackQuery.Data
	log.Printf("收到回調查詢: %s", data)
}

// sendReply 發送回覆訊息
func (h *BotHandlers) sendReply(message *tgbotapi.Message, text string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyToMessageID = message.MessageID

	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("發送訊息時發生錯誤: %v", err)
	}
}

// handleSendToCommand 處理 /sendto 指令 - 發送訊息到指定聊天
func (h *BotHandlers) handleSendToCommand(message *tgbotapi.Message, args string) {
	parts := strings.SplitN(args, " ", 2)
	if len(parts) < 2 {
		h.sendReply(message, "使用方式：/sendto <聊天ID> <訊息內容>\n例如：/sendto -1001234567890 Hello Group!")
		return
	}

	chatIDStr := parts[0]
	messageText := parts[1]

	// 解析聊天ID
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		h.sendReply(message, "無效的聊天ID格式。聊天ID應該是數字。")
		return
	}

	// 發送訊息到指定聊天
	msg := tgbotapi.NewMessage(chatID, messageText)
	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("發送訊息到聊天 %d 失敗: %v", chatID, err)
		h.sendReply(message, fmt.Sprintf("發送失敗：%v", err))
		return
	}

	h.sendReply(message, fmt.Sprintf("✅ 訊息已成功發送到聊天 %d", chatID))
	log.Printf("用戶 %s 發送訊息到聊天 %d: %s", message.From.UserName, chatID, messageText)
}

// handleBroadcastCommand 處理 /broadcast 指令 - 廣播訊息到所有已知聊天
func (h *BotHandlers) handleBroadcastCommand(message *tgbotapi.Message, args string) {
	if args == "" {
		h.sendReply(message, "使用方式：/broadcast <訊息內容>\n例如：/broadcast 重要通知：系統維護中...")
		return
	}

	h.chatMutex.RLock()
	chatIDs := make([]int64, 0, len(h.chatIDs))
	for chatID := range h.chatIDs {
		chatIDs = append(chatIDs, chatID)
	}
	h.chatMutex.RUnlock()

	if len(chatIDs) == 0 {
		h.sendReply(message, "目前沒有已知的聊天。請先讓機器人加入一些群組或頻道。")
		return
	}

	successCount := 0
	failCount := 0

	// 廣播訊息
	for _, chatID := range chatIDs {
		msg := tgbotapi.NewMessage(chatID, args)
		if _, err := h.bot.Send(msg); err != nil {
			log.Printf("廣播到聊天 %d 失敗: %v", chatID, err)
			failCount++
		} else {
			successCount++
		}
	}

	resultText := fmt.Sprintf("📢 廣播完成！\n✅ 成功：%d 個聊天\n❌ 失敗：%d 個聊天",
		successCount, failCount)
	h.sendReply(message, resultText)

	log.Printf("用戶 %s 執行廣播: 成功 %d，失敗 %d", message.From.UserName, successCount, failCount)
}

// SendMessageToChat 公開方法 - 發送訊息到指定聊天（供外部調用）
func (h *BotHandlers) SendMessageToChat(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := h.bot.Send(msg)
	if err != nil {
		log.Printf("發送訊息到聊天 %d 失敗: %v", chatID, err)
	}
	return err
}

// GetKnownChatIDs 取得所有已知的聊天ID
func (h *BotHandlers) GetKnownChatIDs() []int64 {
	h.chatMutex.RLock()
	defer h.chatMutex.RUnlock()

	chatIDs := make([]int64, 0, len(h.chatIDs))
	for chatID := range h.chatIDs {
		chatIDs = append(chatIDs, chatID)
	}
	return chatIDs
}
