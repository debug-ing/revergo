FROM golang:1.23

WORKDIR /app

COPY . .

RUN go build -o main main.go

RUN ls -l /app

ENTRYPOINT ["/app/main"]  