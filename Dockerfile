FROM golang:1.16.3-alpine3.13

WORKDIR /go/src

ENTRYPOINT [ "top" ]