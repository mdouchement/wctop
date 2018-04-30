# build stage
FROM golang:1.10-alpine as build-env
MAINTAINER mdouchement

RUN apk upgrade
RUN apk add --update --no-cache alpine-sdk curl
RUN go get github.com/mjibson/esc
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

ADD . /go/src/github.com/mdouchement/wctop
WORKDIR /go/src/github.com/mdouchement/wctop

RUN dep ensure -v
RUN go generate
RUN go build -o wctop wctop.go

# final stage
FROM alpine:3.7
MAINTAINER mdouchement

COPY --from=build-env /go/src/github.com/mdouchement/wctop/wctop /usr/local/bin/

EXPOSE 8080
CMD ["wctop", "-p", "8080"]
