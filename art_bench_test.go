package art_test

import (
	"runtime"
	"testing"
	"time"

	"github.com/recoilme/art"
)

func forceGC() {
	runtime.GC()
	time.Sleep(time.Millisecond * 500)
}

func BenchmarkSetArt(b *testing.B) {
	keys := seed(b.N, 0)

	b.ResetTimer()
	b.ReportAllocs()
	tree := art.New()
	for n := 0; n < b.N; n++ {
		tree.Set(keys[n], nil)
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

func BenchmarkGetArt(b *testing.B) {
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

func BenchmarkAscendArt(b *testing.B) {
	keys := seed(b.N, 42)

	tree := art.New()
	for n := 0; n < b.N; n++ {
		tree.Set(keys[n], keys[n])
	}

	b.ResetTimer()
	b.ReportAllocs()

	tree.Ascend(nil, func(key, val []byte) bool {
		return true
	})
}

func BenchmarkDescendArt(b *testing.B) {
	keys := seed(b.N, 42)

	tree := art.New()
	for n := 0; n < b.N; n++ {
		tree.Set(keys[n], keys[n])
	}

	b.ResetTimer()
	b.ReportAllocs()

	tree.Descend(nil, func(key, val []byte) bool {
		//fmt.Println(string(key))
		return true
	})
}

func BenchmarkScanArt(b *testing.B) {
	keys := seed(b.N, 42)

	tree := art.New()
	for n := 0; n < b.N; n++ {
		tree.Set(keys[n], keys[n])
	}

	b.ResetTimer()
	b.ReportAllocs()
	tree.Scan(func(key, val []byte) bool {
		return true
	})
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

func BenchmarkGetWordsArt(b *testing.B) {

	keys := loadTestFile("test/words.txt")
	tree := art.New()
	for _, w := range keys {
		tree.Set(w, w)
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
