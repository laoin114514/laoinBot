# syntax=docker/dockerfile:1

FROM registry.cn-hangzhou.aliyuncs.com/library/alpine:3.20

WORKDIR /app

# 使用国内 apk 源
RUN sed -i 's|https://dl-cdn.alpinelinux.org/alpine|https://mirrors.aliyun.com/alpine|g' /etc/apk/repositories \
    && apk add --no-cache ca-certificates tzdata \
    && update-ca-certificates

# 仅拷贝已编译好的可执行文件（不在 Docker 内构建）
COPY laoinBot /app/laoinBot
COPY config /app/config
RUN chmod +x /app/laoinBot

ENV TZ=Asia/Shanghai

ENTRYPOINT ["/app/laoinBot"]
