package art_test

import (
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
