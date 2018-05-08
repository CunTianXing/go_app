package main

import (
	crand "crypto/rand"
	"math/big"
	"math/rand"
	"testing"
)

// 这个基准测试比较了使用Go math / rand和crypto / rand软件包生成伪随机数的性能。
// 使用该math/rand包的随机数生成比使用该包的密码安全随机数生成快得多crypto/rand。
func BenchmarkMathRand(b *testing.B) {
	for n := 0; n < b.N; n++ {
		rand.Int63()
	}
}

func BenchmarkCryptoRand(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := crand.Int(crand.Reader, big.NewInt(27))
		if err != nil {
			panic(err)
		}
	}
}
