FROM golang:1.19.0-alpine3.15 AS builder

RUN apk update
RUN apk upgrade
# RUN apk add bash
RUN apk add git
# RUN apk add certbot

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /app/mochi .

FROM alpine

RUN apk update
RUN apk upgrade
RUN apk add bash
RUN apk add git
RUN apk add certbot

WORKDIR /
COPY --from=builder /app/mochi .
COPY run.sh /run.sh
COPY certbot.sh /certbot.sh
COPY croncert.sh /etc/periodic/weekly/croncert.sh

RUN chmod +x /run.sh
RUN chmod +x /certbot.sh
RUN chmod +x /etc/periodic/weekly/croncert.sh

EXPOSE 80

# tcp
EXPOSE 1883

# websockets
EXPOSE 1882

# dashboard
EXPOSE 8080

# api
EXPOSE 8081

CMD ["/bin/bash","-c","/run.sh"]