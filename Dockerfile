FROM golang:1.19.3 as builder
WORKDIR /app
COPY . ./
RUN go mod download
RUN GOOS=linux go build -a -o telegram-bot .

FROM ubuntu:22.10
RUN adduser --system --group telegram  &&\
    mkdir -p /telegram/data/sqlite &&\
    chown -R telegram /telegram && chgrp -R telegram /telegram &&\
    apt-get update &&\
    apt-get install -y ca-certificates
COPY --from=builder /app/telegram-bot /telegram/
WORKDIR /telegram
USER telegram
CMD ["sh", "-c", "/telegram/telegram-bot -tg-bot-token $TOKEN"]