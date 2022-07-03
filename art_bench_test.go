package art_test

import (
	"encoding/binary"
	"math/rand"
	"runtime"
	"testing"
	"time"

	"github.com/recoilme/art"
)

func randPrintableKey(rnd *rand.Rand, n int) []byte {
	s := make([]byte, n)
	rnd.Read(s)
	for i := 0; i < n; i++ {
		s[i] = 'a' + (s[i] % 26)
	}
	return s
}

func seed(num int, seed int64) [][]byte {

	rng := rand.New(rand.NewSource(seed))
	keys := make([][]byte, num)
	for n := 1; n < num; n++ {
		bin := make([]byte, 8)
		if seed == 0 {
			binary.BigEndian.PutUint64(bin, uint64(n))
		} else {
			binary.BigEndian.PutUint64(bin, rng.Uint64())
		}
		keys[n] = bin
	}
	return keys
}

func forceGC() {
	runtime.GC()
	time.Sleep(time.Millisecond * 500)
}

func BenchmarkSet(b *testing.B) {
	keys := seed(b.N, 0)

	b.ResetTimer()
	b.ReportAllocs()
	tree := art.New()
	for n := 0; n < b.N; n++ {
		tree.Set(keys[n], nil)
	}
}

func BenchmarkGet(b *testing.B) {
	keys := seed(b.N, 0)

	tree := art.New()
	for n := 0; n < b.N; n++ {
		tree.Set(keys[n], nil)
	}

	b.ResetTimer()
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		tree.Get(keys[n])
	}
}
