FROM golang:1.25-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o load-balancer ./cmd/server

EXPOSE 8080

CMD ["./load-balancer"]