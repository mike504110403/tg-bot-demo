# 使用官方 Go 映像作為建構環境
FROM golang:1.21-alpine AS builder

# 設定工作目錄
WORKDIR /app

# 安裝 git（可能需要用於某些 Go 模組）
RUN apk add --no-cache git

# 複製 go mod 檔案
COPY go.mod go.sum ./

# 下載依賴
RUN go mod download

# 複製原始碼
COPY . .

# 建構應用程式
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 使用輕量級的 alpine 映像作為運行環境
FROM alpine:latest

# 安裝 ca-certificates（用於 HTTPS 請求）
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 從建構階段複製執行檔
COPY --from=builder /app/main .

# 設定執行權限
RUN chmod +x ./main

# 暴露埠號（如果使用 webhook 模式）
EXPOSE 8080

# 運行應用程式
CMD ["./main"]
