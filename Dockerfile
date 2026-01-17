FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Сборка
RUN CGO_ENABLED=0 GOOS=linux go build -o app cmd/app/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app .
COPY --from=builder /app/config/config.yaml ./config/config.yaml
COPY --from=builder /app/database ./database

CMD ["./app"]