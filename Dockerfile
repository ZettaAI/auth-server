FROM golang:1.15-buster

WORKDIR /app
COPY . /app

RUN go build
CMD ./auth-server
