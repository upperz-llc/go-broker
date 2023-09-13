FROM golang:1.21.0-alpine3.18 AS builder

RUN apk update
RUN apk upgrade

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /app/mochi cmd/main.go

FROM scratch

WORKDIR /
COPY --from=builder /app/mochi .

ENTRYPOINT [ "/mochi" ]