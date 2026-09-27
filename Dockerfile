
FROM golang:1.25-alpine AS builder

# Добавляем нужные утилиты
RUN apk add --no-cache git ca-certificates

WORKDIR /workspace


RUN apk add --no-cache git \
    && go install github.com/swaggo/swag/cmd/swag@latest
# Сначала копируем только файлы модулей
COPY go.mod go.sum ./


RUN go mod download

# Только после этого копируем весь остальной код
COPY . .

RUN swag init -g cmd/app/main.go

# Компилируем чистый независимый бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o engine ./cmd/app/main.go


# 2. Запуск в минимальном контейнере
FROM alpine:3.19


COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /workspace/engine .

COPY config/config.yaml ./config/config.yaml

EXPOSE 8080
CMD ["./engine"]
