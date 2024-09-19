FROM dockerproxy.com/library/alpine:3.20

# 方便构建镜像时直接把文件拖过去，省去打包步骤
# GOOS=linux GOARCH=amd64  go build -o sgrid_app
# cd client && ./build.sh


# 修改镜像源，并安装所需的依赖
# 也可以进入 容器内部再进行下载
# 为了减少打包体积可以将下面这行去掉
# 其他依赖 nginx=1.26.2-r0
RUN sed -i "s/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g" /etc/apk/repositories \
    && apk add --no-cache ca-certificates \
    openjdk8=8.402.06-r0 \
    nodejs=20.15.1-r0 \
    npm=10.8.0-r0 \
    go=1.22.7-r0

# 设置工作目录
WORKDIR /app

# 将本地代码复制到容器中
COPY ./sgrid_app /app/main
COPY ./dist /app/dist

# 暴露端口
EXPOSE 12111
EXPOSE 14938
EXPOSE 15887

# 启动命令
CMD ["/app/main"]
