# build stage
FROM golang:1.9-alpine as build-env
MAINTAINER mdouchement

RUN apk upgrade
RUN apk add --update --no-cache alpine-sdk
RUN go get github.com/mjibson/esc
RUN go get github.com/Masterminds/glide

ADD . /go/src/github.com/mdouchement/wctop
WORKDIR /go/src/github.com/mdouchement/wctop

RUN glide install
RUN go generate
RUN go build -o wctop wctop.go

# final stage
FROM alpine:3.5
MAINTAINER mdouchement

COPY --from=build-env /go/src/github.com/mdouchement/wctop/wctop /usr/local/bin/

EXPOSE 8080
CMD ["wctop", "-p", "8080"]
