FROM golang:1.13
MAINTAINER https://github.com/andersfylling
WORKDIR /build
COPY cmd/docker /build
RUN go mod download
RUN go test ./...
RUN rm -rf /build
