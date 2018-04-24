package main

import "fmt"
import "testing"

var s string
var x = []byte{1024: 'x'}
var y = []byte{1024: 'y'}
//4
func fc() {
    if string(x) != string(y) {
        s = ("  " + string(x) + string(y))[1:]
    }
}
//2
func fd(){
    if string(x) != string(y) {
        s = string(x) + string(y)
    }
}

func main(){
    fmt.Println(testing.AllocsPerRun(2,fc))//1
    fmt.Println(testing.AllocsPerRun(2,fd))//3
}
