# Telegram Bot Demo - Go 版本

一個使用 Go 語言開發的 Telegram 機器人範例專案。

## 功能特色

- 🤖 完整的 Telegram Bot API 整合
- 💬 支援指令和文字訊息處理
- 📢 支援發送訊息到群組和頻道
- 📡 支援廣播訊息功能
- 🌐 提供 HTTP API 介面
- 🔧 環境變數配置管理
- 🐳 Docker 容器化支援
- 📝 豐富的指令系統
- 🚀 易於擴展的架構

## 快速開始

### 1. 取得 Bot Token

1. 在 Telegram 中找到 [@BotFather](https://t.me/botfather)
2. 發送 `/newbot` 指令
3. 按照指示建立你的機器人
4. 複製取得的 Bot Token

### 2. 設定環境變數

複製範例配置檔案並編輯：

```bash
cp config.env.example .env
```

編輯 `.env` 檔案，填入你的 Bot Token：

```env
TELEGRAM_BOT_TOKEN=your_bot_token_here
DEBUG=true
WEBHOOK_URL=
API_ENABLED=true
API_PORT=8080
```

### 3. 運行機器人

#### 方法一：直接運行

```bash
# 安裝依賴
go mod tidy

# 運行機器人
go run main.go
```

#### 方法二：使用 Docker

```bash
# 建構並運行
docker-compose up -d

# 查看日誌
docker-compose logs -f
```

## 專案結構

```
tg-bot-demo/
├── config/
│   └── config.go          # 配置管理
├── handlers/
│   └── handlers.go        # 訊息處理邏輯
├── api/
│   ├── handlers.go        # API 處理邏輯
│   └── router.go         # API 路由設定
├── main.go               # 主程式進入點
├── go.mod               # Go 模組定義
├── Dockerfile           # Docker 映像配置
├── docker-compose.yml   # Docker Compose 配置
├── config.env.example  # 環境變數範例
└── README.md           # 專案說明
```

## 可用指令

機器人支援以下指令：

- `/start` - 顯示歡迎訊息
- `/help` - 顯示幫助訊息  
- `/info` - 顯示用戶資訊
- `/echo <文字>` - 重複輸入的文字
- `/sendto <聊天ID> <訊息>` - 發送訊息到指定聊天
- `/broadcast <訊息>` - 廣播訊息到所有已知聊天

## 群組功能使用

### 如何取得聊天ID

1. **個人聊天ID**：使用 `/info` 指令查看
2. **群組ID**：將機器人加入群組後，在群組中使用 `/info` 指令
3. **使用 @userinfobot**：向此機器人轉發群組訊息即可取得群組ID

### 群組操作範例

```bash
# 發送訊息到特定群組
/sendto -1001234567890 大家好！這是來自機器人的訊息

# 廣播訊息到所有已知聊天
/broadcast 重要通知：系統將於今晚維護
```

### 程式化發送（供開發者參考）

```go
// 在其他 Go 程式中使用
chatID := int64(-1001234567890) // 群組ID
message := "Hello from Go!"
err := botHandlers.SendMessageToChat(chatID, message)
if err != nil {
    log.Printf("發送失敗: %v", err)
}
```

## HTTP API 使用

機器人提供 RESTful API，讓你可以透過 HTTP 請求發送訊息。

### API 端點

| 方法 | 端點 | 說明 |
|-----|------|------|
| POST | `/api/v1/send` | 發送訊息到指定聊天 |
| POST | `/api/v1/broadcast` | 廣播訊息到所有已知聊天 |
| GET | `/api/v1/chats` | 取得所有已知聊天ID |
| GET | `/health` | 健康檢查 |

### API 使用範例

#### 1. 發送訊息到指定聊天

```bash
curl -X POST http://localhost:8080/api/v1/send \
  -H "Content-Type: application/json" \
  -d '{
    "chat_id": "-1001234567890",
    "message": "Hello from API!"
  }'
```

**回應：**
```json
{
  "success": true,
  "message": "訊息發送成功",
  "chat_id": -1001234567890
}
```

#### 2. 廣播訊息

```bash
curl -X POST http://localhost:8080/api/v1/broadcast \
  -H "Content-Type: application/json" \
  -d '{
    "message": "重要通知：系統維護中"
  }'
```

**回應：**
```json
{
  "success": true,
  "message": "廣播完成",
  "success_count": 3,
  "fail_count": 0
}
```

#### 3. 取得聊天列表

```bash
curl http://localhost:8080/api/v1/chats
```

**回應：**
```json
{
  "success": true,
  "message": "取得聊天列表成功",
  "chat_ids": [-1001234567890, 123456789],
  "count": 2
}
```

#### 4. 健康檢查

```bash
curl http://localhost:8080/health
```

**回應：**
```json
{
  "status": "ok",
  "message": "Telegram Bot API 運行正常"
}
```

### 使用 JavaScript 調用 API

```javascript
// 發送訊息
async function sendMessage(chatId, message) {
  const response = await fetch('http://localhost:8080/api/v1/send', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      chat_id: chatId,
      message: message
    })
  });
  
  const result = await response.json();
  console.log(result);
}

// 廣播訊息
async function broadcast(message) {
  const response = await fetch('http://localhost:8080/api/v1/broadcast', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      message: message
    })
  });
  
  const result = await response.json();
  console.log(result);
}
```

### 使用 Python 調用 API

```python
import requests

# 發送訊息
def send_message(chat_id, message):
    url = "http://localhost:8080/api/v1/send"
    data = {
        "chat_id": str(chat_id),
        "message": message
    }
    response = requests.post(url, json=data)
    return response.json()

# 廣播訊息
def broadcast(message):
    url = "http://localhost:8080/api/v1/broadcast"
    data = {
        "message": message
    }
    response = requests.post(url, json=data)
    return response.json()

# 使用範例
result = send_message(-1001234567890, "Hello from Python!")
print(result)
```

## 開發指令

常用的開發指令：

```bash
# 依賴管理
go mod tidy           # 安裝和整理依賴
go mod download       # 下載依賴

# 開發和建構
go run main.go        # 直接運行
go build -o bot main.go  # 建構二進位檔案
go clean             # 清理建構產物

# 測試
go test ./...         # 運行所有測試
go test -v ./...      # 詳細測試輸出

# 程式碼品質
go fmt ./...          # 格式化程式碼
go vet ./...          # 檢查程式碼

# Docker
docker build -t tg-bot-demo .
docker run --env-file .env tg-bot-demo
```

## 環境變數說明

| 變數名稱 | 說明 | 預設值 | 必要性 |
|---------|------|--------|--------|
| `TELEGRAM_BOT_TOKEN` | Telegram Bot API Token | - | 必要 |
| `DEBUG` | 除錯模式開關 | `true` | 可選 |
| `WEBHOOK_URL` | Webhook URL（如使用 webhook 模式） | - | 可選 |
| `API_ENABLED` | 是否啟用 HTTP API | `true` | 可選 |
| `API_PORT` | HTTP API 伺服器埠號 | `8080` | 可選 |

## 擴展功能

### 新增指令

在 `handlers/handlers.go` 的 `handleCommand` 函數中新增 case：

```go
case "newcommand":
    h.handleNewCommand(message, args)
```

然後實作對應的處理函數：

```go
func (h *BotHandlers) handleNewCommand(message *tgbotapi.Message, args string) {
    // 你的指令邏輯
    h.sendReply(message, "新指令的回應")
}
```

### 新增中介軟體

可以在 `HandleUpdate` 函數中新增中介軟體邏輯，例如：

- 用戶權限檢查
- 訊息記錄
- 速率限制
- 分析統計

### 資料庫整合

建議的資料庫整合方案：

1. **SQLite**：適合小型專案
2. **PostgreSQL**：適合生產環境
3. **Redis**：適合快取和會話管理

## 部署建議

### 開發環境

```bash
go run main.go
```

### 生產環境

```bash
# 使用 Docker Compose
docker-compose up -d

# 或建構二進位檔案
go build -o bot main.go
./bot
```

### 雲端部署

支援部署到：

- Heroku
- AWS ECS
- Google Cloud Run
- Digital Ocean Apps
- 任何支援 Docker 的平台

## 常見問題

### Q: 如何啟用 Webhook 模式？

A: 設定 `WEBHOOK_URL` 環境變數，並在程式中實作 webhook 處理邏輯。

### Q: 如何新增資料庫支援？

A: 在 `config` 包中新增資料庫配置，並建立對應的資料模型和存取層。

### Q: 機器人沒有回應怎麼辦？

A: 檢查：
1. Bot Token 是否正確
2. 機器人是否已啟動
3. 查看日誌輸出
4. 確認網路連線正常

## 貢獻指南

1. Fork 專案
2. 建立功能分支
3. 提交變更
4. 發送 Pull Request

## 授權

此專案採用 MIT 授權。

## 聯絡資訊

如有問題或建議，請開啟 Issue 或聯絡專案維護者。