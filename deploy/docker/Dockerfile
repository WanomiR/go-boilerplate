FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod init studentgit.kata.academy/movie-recommendation-platform/telegram-bot-service \
    && go mod tidy

RUN go build -o ./main ./cmd/app

FROM alpine:latest

COPY --from=builder /app/main /main
COPY --from=builder /app/db/migrations /db/migrations

EXPOSE 8888

CMD ["/main"]