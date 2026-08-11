FROM golang:1.25.0 AS builder

ENV GOOS=linux \
    GOARCH=amd64 \
    GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0

WORKDIR /app

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY . .

RUN go build -o ./knowledge_base_backend ./cmd

FROM alpine:3.19.8

WORKDIR /app

RUN apk update \
    && apk add tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

COPY --chown=0:0 --from=builder /app/knowledge_base_backend /app/knowledge_base_backend
COPY --chown=0:0 --from=builder /app/cmd/knowledge_base_backend_config.ini /app/knowledge_base_backend_config.ini
COPY --chown=0:0 --from=builder /app/migration /app/migration
COPY --chown=0:0 --from=builder /app/config /app/config

CMD ["/app/knowledge_base_backend"]