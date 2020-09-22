FROM golang:1.15.2-buster

WORKDIR /app
COPY . /app

RUN go get -v golang.org/x/oauth2 \
    && go get -v github.com/labstack/echo \
    && go get -v github.com/go-redis/redis