# 使用官方 Golang 镜像作为基础镜像
FROM register.liberx.info/library/golang:1.21.5-alpine as builder

# 设置工作目录
WORKDIR /app

# 将本地代码复制到容器中
COPY . .

# 下载并安装依赖
RUN go mod download && GOOS=linux GOARCH=amd64 go build -o main .

# 构建客户端代码
# RUN cd client && chmod +x build.sh && ./build.sh

# 创建一个更小的镜像用于生产环境
FROM register.liberx.info/library/alpine:3.16

# 安装必要的运行时依赖
RUN apk add --no-cache ca-certificates

# 添加 Java 8 支持
FROM register.liberx.info/library/openjdk:8-jre-alpine as java-builder

FROM register.liberx.info/library/node:16.20-alpine as node-builder

# 设置工作目录
WORKDIR /app

# 将 Java 应用程序复制到容器中
COPY --from=builder /app/main /app/main
COPY --from=builder /app/dist /app/dist
# COPY --from=builder /app/sgrid.yml /app/sgrid.yml


# 暴露端口
EXPOSE 12111
EXPOSE 14938
EXPOSE 15887

# 启动命令
CMD ["/app/main"]
