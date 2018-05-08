package main

import (
    "strconv"
    "testing"
)

func BenchmarkParseBool(b *testing.B) {
    for n := 0; n < b.N; n++ {
        _, err := strconv.ParseBool("true")
        if err != nil {
            panic(err)
        }
    }
}

func BenchmarkParseInt(b *testing.B) {
    for n := 0; n < b.N; n++ {
        _, err := strconv.ParseInt("7182818284",10,64)
        if err != nil {
            panic(err)
        }
    }

}

func BenchmarkParseFloat(b *testing.B) {
    for n := 0; n < b.N; n++ {
        _, err := strconv.ParseFloat("7182818284",64)
        if err != nil {
            panic(err)
        }
    }
}


