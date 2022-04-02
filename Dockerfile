# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.18-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /go-recipes

##
## Deploy
##
FROM gcr.io/distroless/base-debian10

ENV REDIS_SERVER=${REDIS_SERVER}
ENV JWT_SECRET=${JWT_SECRET}
ENV MONGO_URI=${MONGO_URI}
ENV MONGO_DATABASE=${MONGO_DATABASE}
ENV GIN_MODE=${GIN_MODE}


WORKDIR /

COPY --from=build /go-recipes /go-recipes

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/go-recipes"]
