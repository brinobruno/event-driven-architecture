FROM golang:1.23

WORKDIR /app/

RUN apt-get update && apt-get install -y librdkafka-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/walletcore

EXPOSE 7878

ENV PORT=7878

CMD ["./main"]
