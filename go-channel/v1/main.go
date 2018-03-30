package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

//bufio包实现了有缓冲的I/O。
//它包装一个io.Reader或io.Writer接口对象，
//创建另一个也实现了该接口，且同时还提供了缓冲和一些文本I/O的帮助函数的对象。
//func NewScanner(r io.Reader) *Scanner
//NewScanner创建并返回一个从r读取数据的Scanner，默认的分割函数是ScanLines
// Scanner类型提供了方便的读取数据的接口，如从换行符分隔的文本里读取每一行。
// 成功调用的Scan方法会逐步提供文件的token，跳过token之间的字节。token由SplitFunc类型的分割函数指定；默认的分割函数会将输入分割为多个行，并去掉行尾的换行标志。本包预定义的分割函数可以将文件分割为行、字节、unicode码值、空白分隔的word。调用者可以定制自己的分割函数。
// 扫描会在抵达输入流结尾、遇到的第一个I/O错误、token过大不能保存进缓冲时，不可恢复的停止。当扫描停止后，当前读取位置可能会远在最后一个获得的token后面。需要更多对错误管理的控制或token很大，或必须从reader连续扫描的程序，应使用bufio.Reader代替。

// func (s *Scanner) Scan() bool
// Scan方法获取当前位置的token（该token可以通过Bytes或Text方法获得），
//并让Scanner的扫描位置移动到下一个token。当扫描因为抵达输入流结尾或者遇到错误而停止时，本方法会返回false。在Scan方法返回false后，Err方法将返回扫描时遇到的任何错误；除非是io.EOF，此时Err会返回nil。
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lastID := -1
	for scanner.Scan() {
		columns := strings.Split(scanner.Text(), ",")
		id, err := strconv.Atoi(columns[0])
		if err != nil {
			log.Fatalf("ParseInt:%v", err)
		}
		log.Println(columns)
		if id <= lastID {
			log.Fatal("err")
		}
		lastID = id
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("scanner: %v", err)
	}
	log.Println("ok")

}

//tail -n +2 metadata.csv | go run main.go
//tail -n +2 strength_sets.csv | go run main.go
