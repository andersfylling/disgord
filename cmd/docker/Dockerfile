FROM golang:1.13
MAINTAINER https://github.com/andersfylling
WORKDIR /build
COPY . /build
RUN export GO111MODULE=on
RUN go test ./...
RUN rm -rf /build

