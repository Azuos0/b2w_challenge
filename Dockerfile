FROM golang:1.16.3-alpine3.13

ENV CGO_ENABLED 0

WORKDIR /go/src
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .

RUN mkdir ./api

RUN go build -o ../api ./app/main.go

ENTRYPOINT ../api