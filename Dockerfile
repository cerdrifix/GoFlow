FROM golang:1.13.1-alpine3.10

RUN mkdir -p /GoFlow
RUN mkdir -p /certs
RUN mkdir -p /go/src/GoFlow

RUN apk add --no-cache git ca-certificates
RUN go get github.com/jmoiron/sqlx
RUN go get github.com/lib/pq

ADD . /go/src/GoFlow

COPY ./certs /certs

WORKDIR /go/src/GoFlow

RUN go build -o /goFlow .



CMD ["/goFlow"]

EXPOSE 8443