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

func BenchmarkSetHashMap(b *testing.B) {
	keys := seed(b.N, 0)
	m := make(map[string][]byte)
	for _, w := range keys {
		m[string(w)] = w
	}
	b.ResetTimer()
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		m[string(keys[n])] = keys[n]
	}
}

func BenchmarkGetHashMap(b *testing.B) {
	keys := seed(b.N, 0)
	m := make(map[string][]byte)
	for _, w := range keys {
		m[string(w)] = w
	}

	for n := 0; n < b.N; n++ {
		m[string(keys[n])] = keys[n]
	}

	b.ResetTimer()
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		_ = m[string(keys[n])]
	}
}

func BenchmarkGetWords(b *testing.B) {

	keys := loadTestFile("test/words.txt")
	tree := art.New()
	for _, w := range keys {
		_ = tree.Set(w, w)
	}

	b.ResetTimer()
	b.ReportAllocs()
	i := 0
	for n := 0; n < b.N; n++ {
		if n >= len(keys) {
			i = n % len(keys)
		} else {
			i = n
		}
		tree.Get(keys[i])
	}
}

func BenchmarkGetWordsHashMap(b *testing.B) {

	keys := loadTestFile("test/words.txt")
	m := make(map[string][]byte)
	for _, w := range keys {
		m[string(w)] = w
	}

	b.ResetTimer()
	b.ReportAllocs()
	i := 0
	for n := 0; n < b.N; n++ {
		if n >= len(keys) {
			i = n % len(keys)
		} else {
			i = n
		}
		_ = m[string(keys[i])]
	}
}
