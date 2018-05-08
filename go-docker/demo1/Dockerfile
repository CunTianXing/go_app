FROM golang:1.10.1

WORKDIR /go/src/github.com/CunTianXing/go_app/go-docker/demo1/

COPY main.go .

RUN GOOS=linux go build -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=0 /go/src/github.com/CunTianXing/go_app/go-docker/demo1/app .

CMD ["./app"]
