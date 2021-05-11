FROM golang:1.16.3-alpine3.13

ENV CGO_ENABLED 0

WORKDIR /go/src

ENTRYPOINT [ "top" ]