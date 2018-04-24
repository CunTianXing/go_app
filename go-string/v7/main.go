package main

import "fmt"

func main(){
    s := "éक्षिaπ汉字"
    for i, b := range []byte(s){//here, the underlying bytes are not copied.
        fmt.Printf("The byte at index %v: 0x%x \n", i, b)
    }
}
