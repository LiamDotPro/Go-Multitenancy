# deployment file for the go-multitenancy framework
FROM golang:latest

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go build
RUN ./Go-Multitenancy

MAINTAINER Liam Read