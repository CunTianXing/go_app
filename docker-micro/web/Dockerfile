FROM golang:1.9.1

ENV PKG_PATH /go/src/github.com/CunTianXing/go_app/docker-micro

ADD proto $PKG_PATH/proto
ADD google.golang.org $GOPATH/src/google.golang.org
ADD golang.org $GOPATH/src/golang.org
ADD web $PKG_PATH/web
WORKDIR $PKG_PATH/web

RUN go get github.com/gocraft/web github.com/gorilla/context
RUN go get github.com/golang/protobuf/proto

RUN go install

ENTRYPOINT /go/bin/web
EXPOSE 8080
