FROM golang:alpine as builder

RUN apk update
RUN apk add git
RUN apk add make

WORKDIR /usr/src/app

RUN go get -d -v github.com/PuerkitoBio/goquery/
