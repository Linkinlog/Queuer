FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -mod=mod -o queuer -ldflags "-s -w" .

FROM alpine:latest

RUN apk add --update apk-cron && rm -rf /var/cache/apk/*

COPY --from=builder /app/queuer /usr/local/bin/queuer

COPY example-docker.json /usr/local/bin/example.json

RUN echo "* * * * * /usr/local/bin/queuer -f /usr/local/bin/example.json -vv >> /var/log/cron.log 2>&1" > /etc/crontabs/root

CMD ["/usr/sbin/crond", "-f"]
