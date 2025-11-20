FROM golang:1.25-alpine

WORKDIR /cmd

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o cmd .

CMD ["./cmd"]
