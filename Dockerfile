# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.18-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /go-recipes

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/go-recipes"]
