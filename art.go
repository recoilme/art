package art

import (
	"bytes"
	"fmt"
)

type Art struct {
	root *node
}

type node struct {
	key      []byte
	val      []byte
	children []*node
	size     int16
}

func New() *Art {
	return &Art{}
}

func (a *Art) Set(key, val []byte) (replaced bool) {
	if a.root == nil {
		n := &node{}
		n.key = key
		n.val = val
		a.root = n
		return
	}
	replaced = a.root.set(key, val, 0)
	if replaced {
		return
	}
	return
}

func (n *node) set(key, val []byte, depth int) (replaced bool) {
	// leaf node
	if n.size == 0 {
		if bytes.Equal(key[depth:], n.key) {
			n.val = val
			return true
		}
		n.nodeEmptySplit(key[depth:], val)
		return false
	}
	// nodes with children
	if bytes.Equal(key[depth:], n.key) {
		// prefix equal
		n.val = val
		return true
	}
	insert := func() (inserted bool) {
		i, found := n.find(key[depth])
		if found {
			n = n.children[i]
			n.set(key, val, depth)
			// go to child
		} else {
			n.add(key[depth:], val, i)
		}
		return true
	}
	// node without prefix
	if len(n.key) == 0 {
		if insert() {
			return false
		}
	}
	// node with prefix
	cp := commonPrefix(key[depth:], n.key)
	delta := len(n.key) - len(cp)
	if delta <= 0 {
		depth += len(cp)
		if insert() {
			return false
		}
	} else {
		fmt.Println(n.children[0].key)
		// delta>0, prefix: abc, common: a
		// TODO
		//fmt.Println("delta>0", delta, cp, n.key, string(key), key)
	}
	return false
}

func (n *node) grow() {
	switch n.size {
	case 0, 4, 16, 48:
		var newsize int16
		switch n.size {
		case 0:
			newsize = 4
		case 4:
			newsize = 16
		case 16:
			newsize = 48
		case 48:
			newsize = 256
		}
		//newidx := make([]byte, newsize)
		newchild := make([]*node, newsize)
		var i int16
		for i = 0; i < n.size; i++ {
			//newidx[i] = n.idx[i]
			newchild[i] = n.children[i]
		}
		//n.idx = newidx
		n.children = newchild
	}
}

func (n *node) add(key, val []byte, index int16) {
	//	fmt.Println("add", key, n.size)
	// grow
	n.grow()
	var j int16
	for j = n.size; j > index; j-- {
		n.children[j] = n.children[j-1]
	}
	n.children[index] = &node{
		key: key,
		val: val,
	}

	n.size++
}

func (n *node) find(k byte) (index int16, found bool) {
	var low int16
	high := n.size - 1
	for low <= high {
		mid := low + ((high+1)-low)/2
		if k >= n.children[mid].key[0] {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	if low > 0 && n.children[low-1].key[0] == k {
		index = low - 1
		found = true
	} else {
		index = low
		found = false
	}
	return index, found
}

func commonPrefix(key, newKey []byte) []byte {
	i := 0
	limit := min(len(key), len(newKey))
	for ; i < limit; i++ {
		if key[i] != newKey[i] {
			return key[:i]
		}
	}
	return key[:i]
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func (n *node) nodeEmptySplit(key, val []byte) {
	//fmt.Println("nodeEmptySplit", key)
	cp := commonPrefix(key, n.key)
	depth := len(cp)
	// save old vals
	oldKey := n.key
	oldVal := n.val
	// modify to node with childs
	n.key = cp
	n.val = nil
	n.children = make([]*node, 4)
	// add childs
	n.add(oldKey[depth:], oldVal, 0)
	ind, _ := n.find(key[depth])
	n.add(key[depth:], val, ind)
}
