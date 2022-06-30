package art_test

import (
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/recoilme/art"
)

func Test1(t *testing.T) {
	art := art.New()
	item := []byte("1")
	replaced := art.Set(item, item)
	if replaced {
		t.Fatal("expected false")
	}
	replaced = art.Set(item, []byte("2"))
	if !replaced {
		t.Fatal("expected true")
	}
}

func Test2(t *testing.T) {
	art := art.New()
	item1 := []byte("1")
	item2 := []byte("2")
	replaced := art.Set(item2, item2)
	if replaced {
		t.Fatal("expected false")
	}
	replaced = art.Set(item1, item1)
	if replaced {
		t.Fatal("expected false")
	}
}

func Test3(t *testing.T) {
	art := art.New()
	item1 := []byte("api1")
	item2 := []byte("api2")
	item3 := []byte("a")
	replaced := art.Set(item2, item2)
	if replaced {
		t.Fatal("expected false")
	}
	replaced = art.Set(item1, item1)
	if replaced {
		t.Fatal("expected false")
	}
	art.Set(item3, item3)
	if replaced {
		t.Fatal("expected false")
	}
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
	for n := b.N - 1; n >= 0; n-- {
		tree.Set(strs[n], strs[n])
	}
	fmt.Println(fmt.Sprintf("%+v", tree))
}

func Test4(t *testing.T) {
	N := 3
	strs := make([][]byte, N)

	for n := 0; n < N; n++ {
		bin := make([]byte, 8)
		binary.BigEndian.PutUint64(bin, uint64(n))
		strs[n] = bin
	}

	tree := art.New()
	for n := N - 1; n >= 0; n-- {
		tree.Set(strs[n], strs[n])

		if n == 3 {
			tree.Set(strs[n], strs[n])
			break //fmt.Println(fmt.Sprintf("%+v", tree))
		}
	}
	fmt.Println(fmt.Sprintf("%+v", tree))
}
