FROM golang:1.22.10-alpine3.21

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o main ./cmd

CMD ["./main"]

