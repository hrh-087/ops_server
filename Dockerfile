FROM golang:alpine AS builder

WORKDIR /data/ops
LABEL stage=gobuilder
COPY . .

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct
ENV GO111MODULE on

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk update --no-cache && apk add --no-cache tzdata

ADD go.mod .
ADD go.sum .
ADD script .

RUN go env
RUN go mod tidy
RUN go build -ldflags="-s -w" -o /data/ops/server .

FROM alpine:latest

LABEL MAINTAINER="dc"

# 设置时区
ENV TZ=Asia/Shanghai
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk update && apk add --no-cache tzdata openntpd && ln -sf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
WORKDIR /data/ops

COPY --from=0 /data/ops/server /data/ops/server
COPY --from=0 /data/ops/game_script /data/ops/game_script
COPY --from=0 /data/ops/config.yaml /data/ops/etc/config.yaml
#COPY --from=0 /data/ops/resource /data/ops/resource

EXPOSE 8000
ENTRYPOINT ./server -c /data/ops/etc/config.yaml