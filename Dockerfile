# syntax=docker/dockerfile:1

FROM golang:1.23

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download
RUN go mod tidy

COPY . .

RUN go build -o main ./cmd/app

EXPOSE 8080

CMD ["./main"]