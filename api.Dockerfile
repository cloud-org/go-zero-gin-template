# 两阶段构建 master
FROM golang:1.16.2-alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOOS linux
ENV GO111MODULE on
ENV GOPROXY https://goproxy.cn,direct

WORKDIR /build/zero

COPY . .
COPY api/etc /app/etc
RUN go build -ldflags="-s -w" -o /app/api api/api.go


FROM alpine:3.15

RUN echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.15/main/" > /etc/apk/repositories \
&& apk update --no-cache && apk add --no-cache ca-certificates tzdata
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/api /app/api
COPY --from=builder /app/etc /app/etc

EXPOSE 9097

ENTRYPOINT ["./api"]

