FROM golang:1.20-alpine as builder

# 设置工作目录
WORKDIR /cmd

# 复制go.mod和go.sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制项目文件
COPY . .

# 构建Go程序
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# 使用轻量级的Alpine Linux作为基础镜像
FROM alpine:latest

# 安装必要的依赖
RUN apk --no-cache add ca-certificates

# 设置工作目录
WORKDIR /root/

# 从builder阶段复制构建好的二进制文件
COPY --from=builder /cmd/main .

# 暴露端口（根据你的程序实际情况修改）
EXPOSE 8080
