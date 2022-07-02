package art_test

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/recoilme/art"
)

func Test1(t *testing.T) {
	art := art.New()
	item := []byte("1")
	replaced := art.Set(item, nil)
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
	replaced = art.Set(item3, item3)
	if replaced {
		t.Fatal("expected false")
	}
	//art.Print()
}

func Test4(t *testing.T) {
	N := 5
	strs := make([][]byte, N)

	for n := 0; n < N; n++ {
		bin := make([]byte, 8)
		binary.BigEndian.PutUint64(bin, uint64(n))
		strs[n] = bin
	}

	tree := art.New()
	for n := N - 1; n >= 0; n-- {
		if n == 0 {
			tree.Set(strs[n], strs[n])
		}
		tree.Set(strs[n], strs[n])
	}
	t.Log(tree.String())
}

func Test5(t *testing.T) {
	N := 23
	strs := make([][]byte, N)

	for n := 0; n < N; n++ {
		bin := make([]byte, 8)
		binary.BigEndian.PutUint64(bin, uint64(n))
		strs[n] = bin
	}

	tree := art.New()
	for n := 0; n < N; n++ {
		tree.Set(strs[n], strs[n])
	}
	for n := 0; n < N; n++ {
		val := tree.Get(strs[n])
		if !bytes.Equal(strs[n], val) {
			t.Fail()
		}
	}
}
