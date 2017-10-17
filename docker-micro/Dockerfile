FROM golang:1.9.1

ARG SRV_NAME

ENV PKG_PATH /go/src/github.com/CunTianXing/go_app/docker-micro

ADD proto $PKG_PATH/proto
ADD shared $PKG_PATH/shared
ADD google.golang.org $GOPATH/src/google.golang.org
ADD golang.org $GOPATH/src/golang.org
ADD gopkg.in $GOPATH/src/gopkg.in
ADD $SRV_NAME $PKG_PATH/$SRV_NAME
WORKDIR $PKG_PATH/$SRV_NAME

RUN go get github.com/golang/protobuf/proto

RUN go install
EXPOSE 8080
