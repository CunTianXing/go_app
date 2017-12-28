package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/proxy"
)

//Socks5Client ...
func Socks5Client(addr string, auth ...*proxy.Auth) (client *http.Client, err error) {
	dialer, err := proxy.SOCKS5("tcp", addr, nil, &net.Dialer{Timeout: 30 * time.Second, KeepAlive: 30 * time.Second})
	if err != nil {
		return
	}
	transport := &http.Transport{Proxy: nil, Dial: dialer.Dial, TLSHandshakeTimeout: 10 * time.Second}
	client = &http.Client{Transport: transport}
	return
}
func main() {
	client, err := Socks5Client("172.20.14.83:8080")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	resp, err := client.Get("http://mengqi.info")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(b))
	}
}
