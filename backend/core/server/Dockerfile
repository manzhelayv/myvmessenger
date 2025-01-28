FROM golang:1.22 as builder

ENV CGO_ENABLED=0

COPY . /app
WORKDIR /app

RUN go build -a -installsuffix cgo -o main /app/*.go