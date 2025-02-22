# 构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 创建 go.mod 文件
RUN go mod init mrd

# 复制源代码
COPY . .

# 添加依赖到 go.mod
RUN go get github.com/gin-gonic/gin && \
    go get github.com/joho/godotenv && \
    go mod tidy

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 运行阶段
FROM alpine:latest

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/.env .

# 创建存储目录
RUN mkdir -p storage

EXPOSE 8080

CMD ["./main"] 