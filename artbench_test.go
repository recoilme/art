package art_test

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/recoilme/art"
)

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
