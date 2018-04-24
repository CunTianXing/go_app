package main

import "fmt"
//字符串可以用作特殊情况下的字节片。

func main(){
    hello := []byte("Hello")
    world := "world"

    helloWorld1 := append(hello,[]byte(world)...)
    fmt.Println(string(helloWorld1))

    helloWorld2 := append(hello,world...)
    fmt.Println(string(helloWorld2))

    helloWorld3 := make([]byte,len(hello)+len(world))
    copy(helloWorld3,hello)
    //copy(helloWorld3[len(hello):],[]byte(world))
    copy(helloWorld3[len(hello):],world)
    fmt.Println(string(helloWorld3))
}
