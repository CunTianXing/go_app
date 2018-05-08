package main

import (
	crand "crypto/rand"
	"math/rand"
	"testing"
)

// 该基准比较了基于Go math / rand和crypto / rand的16个字符，均匀分布的随机字符串生成的性能。
// 使用math/rand包的随机字符串生成比使用该crypto/rand包的密码安全的随机字符串生成要快。
const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz+/"

func randomBytes(n int) []byte {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return bytes
}

func cryptoRandomBytes(n int) []byte {
	bytes := make([]byte, n)
	_, err := crand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return bytes
}

func randomString(bytes []byte) string {
	for i, b := range bytes {
		bytes[i] = letters[b%64]
	}
	return string(bytes)
}

func BenchmarkMathRandString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		randomString(randomBytes(16))
	}
}

func BenchmarkCryptoRandString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		randomString(cryptoRandomBytes(16))
	}
}
