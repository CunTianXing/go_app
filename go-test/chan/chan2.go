package main

import "fmt"

func get_notification(user string) chan string {
	notifications := make(chan string)
	go func() {
		notifications <- fmt.Sprintf("Hi %s, welcome to weibo.com!", user)
	}()
	return notifications
}

func main() {
	jack := get_notification("xingcuntian")
	gary := get_notification("gary")

	fmt.Println(<-jack)
	fmt.Println(<-gary)
}
