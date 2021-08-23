FROM golang:latest AS builder

ENV GOPATH=/go
#ENV TZ=Asia/Shanghai
WORKDIR $GOPATH/src/vault
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o vault main.go

FROM alpine:latest
WORKDIR /server/
COPY --from=builder /go/src/vault/vault .
EXPOSE 8081
# 设置时区为上海
RUN apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata \
    && apk add --no-cache bash
CMD ["./vault"]