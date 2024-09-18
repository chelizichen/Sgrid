# 事先打包好再丢过去

FROM dockerproxy.com/library/alpine:3.20


# 修改镜像源，并安装所需的依赖
RUN sed -i "s/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g" /etc/apk/repositories \
    && apk add --no-cache ca-certificates \
    openjdk8=8.402.06-r0 \
    nodejs=20.15.1-r0 \
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
