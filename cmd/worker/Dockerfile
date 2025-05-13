
FROM golang:1.23.4-alpine3.21 AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o worker ./cmd/worker

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache sqlite ca-certificates

COPY --from=builder /app/worker .

RUN mkdir -p /data

CMD ["./worker", "--db", "/data/trades.db", "--poll", "100ms"]