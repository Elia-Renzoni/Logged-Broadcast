# image build
FROM golang:1.25

WORKDIR /cmd

COPY go.mod go.sum

CMD ['/cmd']

