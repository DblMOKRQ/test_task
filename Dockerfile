FROM golang:1.25.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /app/server ./cmd/

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .

COPY ./migrations ./migrations

EXPOSE 8080

CMD ["/app/server"]