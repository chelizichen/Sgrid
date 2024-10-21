FROM dockerproxy.net/library/alpine:3.20

RUN sed -i "s/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g" /etc/apk/repositories \
    && apk add --no-cache ca-certificates \
    openjdk8=8.402.06-r0 \
    nodejs=20.15.1-r0 \
    npm=10.8.0-r0 \
    go=1.22.8-r0

WORKDIR /app

COPY ./sgrid_app /app/main
COPY ./dist /app/dist

EXPOSE 12111
EXPOSE 14938
EXPOSE 15887

CMD ["/app/main"]
