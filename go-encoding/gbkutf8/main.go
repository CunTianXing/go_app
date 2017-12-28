package main

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

//GbkToUtf8 ...
func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// Utf8ToGbk ...
func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return d, nil
}

//Golang 中的 UTF-8 与 GBK 编码转换...
func main() {
	s := "GBK 与 UTF-8 编码转换测试"
	gbk, err := Utf8ToGbk([]byte(s))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(gbk))

	utf8, err := GbkToUtf8(gbk)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(utf8))
}
