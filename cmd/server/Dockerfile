FROM golang:1.23.4-alpine3.21 AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o server ./cmd/server

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache sqlite ca-certificates

COPY --from=builder /app/server .

RUN mkdir -p /data

EXPOSE 8080

CMD ["./server", "--db", "/data/trades.db", "--listen", "8080"]