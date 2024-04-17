package tools

import (
	"bytes"
	"compress/gzip"
	"io"
	"math/rand"
	"testing"
	"time"
)

func BenchmarkGenerateCrypto(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateCrypto(16)
	}
}

func BenchmarkGenerateShortKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateShortKey()
	}
}

func BenchmarkGenerateCookie(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		n := r.Intn(100)
		b.StartTimer()
		GenerateCookie(n)
	}
}

func BenchmarkUnpacking(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()

		var byt bytes.Buffer
		w := gzip.NewWriter(&byt)
		w.Write([]byte(RandomString()))
		w.Close()

		b.StartTimer()

		Unpacking(io.NopCloser(bytes.NewReader(byt.Bytes())))
	}
}

func RandomString() string {
	mystr := ""
	symbols := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := 2 + r.Intn(10)

	for i := 0; i <= n; i++ {
		index := r.Intn(62)
		mystr += string(symbols[index])
	}

	return mystr
}
