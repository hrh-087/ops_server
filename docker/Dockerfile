FROM golang:alpine as builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct
ENV GO111MODULE on

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /data/ops

COPY . .
ADD go.mod .
ADD go.sum .

RUN go env
RUN go mod tidy
RUN go build -ldflags="-s -w" -o /data/ops/server .

FROM alpine:latest

LABEL MAINTAINER="dc"

# 设置时区
ENV TZ=Asia/Shanghai
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk update && apk get --no-cache tzdata openntpd && ln -sf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
WORKDIR /data/ops

COPY --from=builder /data/ops/server /data/ops/server
COPY --from=builder /data/ops/resource /data/ops/resource
COPY --from=builder /data/ops/config.yaml /data/ops/config.yaml

EXPOSE 8000
ENTRYPOINT ./server -c config.yaml