package art_test

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/recoilme/art"
)

func initBin(N int) [][]byte {
	bytes := make([][]byte, N)
	for n := 0; n < N; n++ {
		bytes[n] = make([]byte, 8)
		binary.BigEndian.PutUint64(bytes[n], uint64(n))
	}
	return bytes
}

func BenchmarkSet(b *testing.B) {
	strs := make([][]byte, b.N)

	for n := 0; n < b.N; n++ {
		bin := make([]byte, 8)
		binary.BigEndian.PutUint64(bin, uint64(n))
		strs[n] = bin
	}

	b.ResetTimer()
	b.ReportAllocs()
	tree := art.New()
	for n := b.N - 1; n > 0; n-- {
		tree.Set(strs[n], nil)
	}
}

func BenchmarkGet(b *testing.B) {
	strs := make([][]byte, b.N)

	for n := 0; n < b.N; n++ {
		bin := make([]byte, 8)
		binary.BigEndian.PutUint64(bin, uint64(n))
		strs[n] = bin
	}

	tree := art.New()
	for n := b.N - 1; n > 0; n-- {
		tree.Set(strs[n], strs[n])
	}
	b.ResetTimer()
	b.ReportAllocs()
	for n := b.N - 1; n > 0; n-- {
		val := tree.Get(strs[n])
		if !bytes.Equal(val, strs[n]) {
			b.Fail()
		}
	}
}

func BenchmarkArt4Get(b *testing.B) {
	strs := initBin(b.N)

	tree := art.New()
	for n := 0; n < b.N; n++ {
		tree.Set(strs[n], nil)
	}

	b.ResetTimer()
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		_ = tree.Get(strs[n])

	}
}
