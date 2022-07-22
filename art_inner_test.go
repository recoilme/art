package art

import (
	"fmt"
	"sort"
	"testing"
)

func TestPrefix2(t *testing.T) {
	dataSet := []struct {
		keyPrefix string
		keys      []string
		expected  []string
	}{
		{
			"long.api.url.v",
			[]string{"long.api.url.v1.foo", "long.api.url.v1.bar", "long.api.url.v2.foo"},
			[]string{"long.api.url.v1.foo", "long.api.url.v1.bar"},
		},
		{
			"this:key:has",
			[]string{
				"this:key:has:a:long:prefix:3",
				"this:key:has:a:long:common:prefix:2",
				"this:key:has:a:long:prefix:33",
				"this:key:has:a:long:prefix:333",
				"this:key:has:a:long:common:prefix:1",
			},
			[]string{
				"this:key:has:a:long:prefix:3",
				"this:key:has:a:long:common:prefix:2",
				"this:key:has:a:long:common:prefix:1",
			},
		},
	}

	for _, d := range dataSet {
		tree := New()
		for _, k := range d.keys {
			tree.Set([]byte(k), []byte(k))

		}

		sort.Sort(sort.Reverse(sort.StringSlice(d.expected)))
		n := tree.root
		fmt.Println(n.StringKeys(0, true))
		//root := n
		pivot := []byte(d.keyPrefix)
		n, depth := n.get(pivot, 0, false)
		//fmt.Println(":" + string(n.key))
		if n != nil {
			cs := commonSuffix(pivot[:depth], n.key)
			depth -= len(cs)
			pref := pivot[:depth]
			parents := make([][]byte, 0, 8)
			n.parents2(func(key, val []byte) bool {
				fmt.Println(string(key), ":", string(val))
				parents = append(parents, key)
				return true
			}, pref, len(pref))

		}

		/*
			tree.Descend([]byte(d.keyPrefix), leafFilter)
			for i := range d.expected {
				_ = i
				//if d.keyPrefix == "api" {
				if !bytes.Equal([]byte(d.expected[i]), []byte(actual[i])) {
					t.Error("Bad news:", d.keyPrefix, actual, d.expected)
					tree.Descend([]byte(d.keyPrefix), leafFilter)
					t.Fatal(tree.StringKeys(true))
				}
				//}
			}*/
	}

}
