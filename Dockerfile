FROM golang:latest

RUN go get github.com/thinkong/myclone
RUN go install github.com/thinkong/myclone/main
ADD ./templates ./templates
ENTRYPOINT /go/bin/main

EXPOSE 8080