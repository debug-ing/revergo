FROM golang:1.23 AS builder

WORKDIR /app

COPY . .

RUN go build -o main main.go

ENTRYPOINT ["/app/main"]  