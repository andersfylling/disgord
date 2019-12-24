FROM golang:1.13
MAINTAINER https://github.com/andersfylling
WORKDIR /build
COPY cmd/script /build
RUN go mod download
RUN rm -rf /build
