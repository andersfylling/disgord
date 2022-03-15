FROM golang:1.17
MAINTAINER https://github.com/andersfylling
WORKDIR /build
COPY cmd/dockermain.go /build
RUN go mod init mybot
RUN go test .
RUN rm -rf /build
