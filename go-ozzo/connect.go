package main

import (
    "fmt"
    "log"
    _ "github.com/go-sql-driver/mysql"
    "github.com/go-ozzo/ozzo-dbx"
)

func main() {
    db, err := dbx.Open("mysql","xingcuntian:xingcuntian@tcp://127.0.0.1:3306/xingcuntian")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(db)
}
