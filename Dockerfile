FROM golang:latest AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o garantex-monitor ./cmd

FROM alpine:latest AS runner

WORKDIR /app

COPY --from=builder /app/garantex-monitor .
COPY --from=builder /app/.env .
COPY --from=builder /app/migrations ./migrations

CMD ["/app/garantex-monitor"]