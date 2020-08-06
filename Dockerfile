FROM golang:alpine AS builder
ENV BIN_PATH="/go/bin/bot"
WORKDIR /go/src/bot
RUN apk update && apk add --no-cache git ca-certificates \
 && update-ca-certificates \
 && go get github.com/go-telegram-bot-api/telegram-bot-api
COPY . .
RUN go version
RUN GOOS=linux GOARCH=386 go build -ldflags="-w -s" -o ${BIN_PATH} \
 && chmod +x /go/bin/bot

FROM scratch 
ENV APP_OUTPUT_FILE ${APP_OUTPUT_FILE:-"/buffer/app.output"}
ENV CMD_OUTPUT_FILE ${CMD_OUTPUT_FILE:-"/buffer/cmd.output"}
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/bot /go/bin/bot
ENTRYPOINT ["/go/bin/bot"]
