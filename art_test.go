package art_test

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"math/rand"
	"os"
	"sort"
	"testing"

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
	for n := 0; n < num; n++ {
		bin := make([]byte, 8)
		if seed == 0 {
			binary.BigEndian.PutUint64(bin, uint64(n))
		} else {
			if seed == 42 {
				bin = randPrintableKey(rng, 8)
			} else {
				binary.BigEndian.PutUint64(bin, rng.Uint64())
			}
		}
		keys[n] = bin
	}
	return keys
}

func Test1(t *testing.T) {
	art := art.New()
	item := []byte("1")
	art.Set(item, nil)
	art.Set(item, []byte("2"))
	if !bytes.Equal([]byte("2"), art.Get(item)) {
		t.Fail()
	}
}

func Test2(t *testing.T) {
	art := art.New()
	item1 := []byte("1")
	item2 := []byte("2")
	art.Set(item2, item2)
	art.Set(item1, item1)
	if !bytes.Equal(item1, art.Get(item1)) {
		t.Fail()
	}
	if !bytes.Equal(item2, art.Get(item2)) {
		t.Fail()
	}
}

func Test3(t *testing.T) {
	art := art.New()
	item1 := []byte("api1")
	item2 := []byte("api2")
	item3 := []byte("a")
	art.Set(item2, item2)
	art.Set(item1, item1)
	art.Set(item3, item3)
	if !bytes.Equal(item1, art.Get(item1)) {
		t.Fail()
	}
	if !bytes.Equal(item2, art.Get(item2)) {
		t.Fail()
	}
	if !bytes.Equal(item3, art.Get(item3)) {
		t.Fail()
	}
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
		tree.Set(strs[n], strs[n])
	}
	for n := N - 1; n >= 0; n-- {
		if !bytes.Equal(tree.Get(strs[n]), strs[n]) {
			t.Fail()
		}
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
		t.Error("unexpected result")
	}

	if res := tree.Get([]byte("yolo")); !bytes.Equal(res, earth) {
		t.Error("unexpected result")
	}

	if res := tree.Get([]byte("yoli")); !bytes.Equal(res, earth) {
		t.Error("unexpected result")
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

func TestRandomInt(t *testing.T) {
	items := seed(1_000, 43)
	tree := art.New()
	for _, item := range items {
		tree.Set(item, item)
	}
	for _, item := range items {
		if !bytes.Equal(item, tree.Get(item)) {
			t.Fail()
		}
	}
	//t.Log(tree.StringKeys(true))
}

func TestRandomBin(t *testing.T) {
	items := seed(1_000, 42)
	tree := art.New()
	for _, item := range items {
		tree.Set(item, item)
	}
	for _, item := range items {
		if !bytes.Equal(item, tree.Get(item)) {
			t.Fail()
		}
	}
	//t.Log(tree.StringKeys(true))
}

func TestDelete(t *testing.T) {
	art := art.New()
	item := []byte("1")
	item2 := []byte("2")
	art.Set(item, item)
	art.Delete(item)
	if art.Get(item) != nil {
		t.Fail()
	}
	art.Set(item, item)
	art.Set(item2, item2)
	art.Delete(item2)
	if art.Get(item2) != nil {
		t.Fail()
	}
	art.Delete(item)
	if art.Get(item) != nil {
		t.Fail()
	}
	t.Log(art.String())
}

func TestClean(t *testing.T) {
	items := seed(100000, 42)
	tree := art.New()
	for _, item := range items {
		tree.Set(item, item)
	}

	for _, item := range items {
		tree.Delete(item)
	}

	for _, item := range items {
		if tree.Get(item) != nil {
			t.Fail()
		}
	}

	t.Log(tree.String())
}

func TestScan(t *testing.T) {
	tree := art.New()
	items := seed(50000, 42)
	items = append(items,
		[]byte("f1"), []byte("n1"), []byte("n11"))
	for _, item := range items {
		tree.Set(item, item)
	}
	//t.Log(tree.StringKeys(true))
	newitems := make([][]byte, 0, len(items))
	tree.Scan(func(key, val []byte) bool {
		if !bytes.Equal(key, val) {
			t.Fatal("not equal", key, val)
		}
		newitems = append(newitems, []byte(string(key)))
		return true
	})
	sort.Slice(items, func(i, j int) bool {
		return bytes.Compare(items[i], items[j]) <= 0
	})
	for i := range items {
		if !bytes.Equal(items[i], newitems[i]) {
			t.Fatal("not equal", items[i], newitems[i])
		}
	}
}

func TestAscend(t *testing.T) {
	tree := art.New()

	earth := []byte("earth")

	tree.Set([]byte("hello"), []byte("world"))
	tree.Set([]byte("yo"), earth)
	tree.Set([]byte("yolo"), earth)
	tree.Set([]byte("yol"), earth)
	tree.Set([]byte("yoli"), earth)
	tree.Set([]byte("yopo"), earth)
	tree.Set([]byte("he"), earth)
	tree.Set([]byte("haa"), earth)
	tree.Set([]byte("hab"), earth)
	//t.Log(tree.StringKeys(true))
	var last []byte
	tree.Ascend([]byte("yo"), func(key, val []byte) bool {
		if !bytes.HasPrefix(key, []byte("yo")) {
			//t.Fatal()
		}
		if bytes.Compare(key, last) < 0 {
			//t.Fatal("out of order")
		}
		last = key
		return true
	})
}

func TestDescend(t *testing.T) {
	tree := art.New()

	tree.Set([]byte("hello"), []byte("hello"))
	tree.Set([]byte("yo"), []byte("yo"))
	tree.Set([]byte("yolo"), []byte("yolo"))
	tree.Set([]byte("yol"), []byte("yol"))
	tree.Set([]byte("yoli"), []byte("yoli"))
	tree.Set([]byte("yopo"), []byte("yopo"))
	tree.Set([]byte("he"), []byte("he"))
	tree.Set([]byte("haa"), []byte("haa"))
	tree.Set([]byte("hab"), []byte("hab"))
	//t.Log(tree.StringKeys(true))
	var last []byte
	pivot := []byte("h")
	tree.Descend(pivot, func(key, val []byte) bool {
		if last == nil {
			last = key
		}
		if !bytes.HasPrefix(key, pivot) {
			t.Fatal()
		}
		if bytes.Compare(key, last) > 0 {
			t.Fatal("out of order")
		}
		//t.Log(string(key), string(val))
		last = key
		return true
	})

	N := 10
	keys := seed(N, 42)

	tree = art.New()
	for n := 0; n < N; n++ {
		tree.Set(keys[n], keys[n])
	}

	tree.Descend(nil, func(key, val []byte) bool {
		//fmt.Println(string(key))
		return true
	})
}
