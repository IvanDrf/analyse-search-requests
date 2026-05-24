FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o search_service ./cmd/main.go

FROM alpine:3.18

WORKDIR /app
COPY --from=builder /app/search_service .
COPY --from=builder /app/config/ ./config/

EXPOSE 8080
CMD ["./search_service", "--config=./config/config.yaml"]