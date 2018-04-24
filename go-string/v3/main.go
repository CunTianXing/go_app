package main

import (
    "fmt"
    "bytes"
    "unicode/utf8"
)

func Runes2Bytes(rs []rune) []byte{
    n := 0
    fmt.Println(rs)
    for _, r := range rs {
        n +=utf8.RuneLen(r)
    }

    n, bs := 0, make([]byte, n)
    for _, r := range rs {
        fmt.Println(string(r))
        n += utf8.EncodeRune(bs[n:],r)
    }
    return bs
}

func main(){
    s := "xingcuntian"
    bs := []byte(s) // string -> []byte
    s = string(bs) // []byte -> string
    rs := []rune(s) // string -> []rune
    s = string(rs)  // []rune -> string
    rs = bytes.Runes(bs) // []byte -> []rune
    bs = Runes2Bytes(rs) // []rune -> []byte
    fmt.Println(rs)
    fmt.Println(bs)
}
