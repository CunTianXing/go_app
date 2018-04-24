package main

import "fmt"

func main(){
    s := "éक्षिaπ汉字"
    for i:=0; i<len(s); i++ {
        fmt.Printf("The byte at index %v: 0x%x \n",i , s[i])
    }
}
