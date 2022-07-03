package art_test

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"os"
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

func TestTreeInsert2AndSearch(t *testing.T) {
	tree := art.New()

	earth := []byte("earth")

	tree.Set([]byte("hello"), []byte("world"))
	tree.Set([]byte("yo"), earth)
	tree.Set([]byte("yolo"), earth)
	tree.Set([]byte("yol"), earth)
	tree.Set([]byte("yoli"), earth)
	tree.Set([]byte("yopo"), earth)
	t.Log(tree.StringKeys(true))

	if res := tree.Get([]byte("yo")); !bytes.Equal(res, earth) {
		t.Error("unexpected search result")
	}

	if res := tree.Get([]byte("yolo")); !bytes.Equal(res, earth) {
		t.Error("unexpected search result")
	}

	if res := tree.Get([]byte("yoli")); !bytes.Equal(res, earth) {
		t.Error("unexpected search result")
	}
}

func TestStringKeys(t *testing.T) {
	tree := art.New()

	tree.Set([]byte("http://example.com/tag/10"), []byte("a"))
	tree.Set([]byte("http://example.com/tag/20"), []byte("b"))
	tree.Set([]byte("http://some.com"), []byte("c"))

	t.Log(tree.StringKeys(true))
}

func loadTestFile(path string) [][]byte {
	file, err := os.Open(path)
	if err != nil {
		panic("Couldn't open " + path)
	}
	defer file.Close()

	var words [][]byte
	reader := bufio.NewReader(file)
	for {
		if line, err := reader.ReadBytes(byte('\n')); err != nil {
			break
		} else {
			if len(line) > 0 {
				words = append(words, line[:len(line)-1])
			}
		}
	}
	return words
}

func TestWords(t *testing.T) {
	worlds := loadTestFile("test/words.txt")
	tree := art.New()
	for _, w := range worlds {
		tree.Set(w, w)
	}
	for _, w := range worlds {
		if !bytes.Equal(w, tree.Get(w)) {
			t.Fail()
		}
	}
}
