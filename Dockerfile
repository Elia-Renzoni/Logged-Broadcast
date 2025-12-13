FROM golang:1.25

WORKDIR /app

RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o node ./cmd

CMD ["./cmd/main"]

