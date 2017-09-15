FROM golang:alpine as builder

WORKDIR /go/src/

RUN apk update
RUN apk add git
RUN apk add make

RUN make build
