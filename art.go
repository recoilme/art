package art

import (
	"bytes"
	"fmt"
	"math/bits"
	"strings"
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
	/*newKey := make([]byte, len(key))
	copy(newKey, key)
	newVal := make([]byte, len(val))
	copy(newVal, val)*/
	newKey := key
	newVal := val
	if a.root == nil {
		n := &node{}
		n.key = newKey
		n.val = newVal
		a.root = n
		return
	}
	replaced = a.root.set(newKey, newVal, 0)
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
		// rebase tree
		oldnode := &node{}
		*oldnode = *n
		oldnode.key = oldnode.key[len(cp):]
		//oldnode.Print(0)
		newnode := &node{
			key:      cp,
			val:      nil,
			children: make([]*node, 4),
		}
		newnode.children[0] = oldnode
		newnode.size = 1
		if bytes.Equal(key[depth:], cp) {
			newnode.val = val
			*n = *newnode
		} else {
			*n = *newnode
			n.set(key, val, depth)
		}
		return false
	}
	return false
}

func (n *node) grow() {
	switch n.size {
	case 4, 16, 48:
		var newsize int
		switch n.size {
		//case 0:
		//	newsize = 4
		case 4:
			newsize = 16
		case 16:
			newsize = 48
		case 48:
			newsize = 256
		}
		newchild := make([]*node, newsize)
		var i int16
		for i = 0; i < n.size; i++ {
			newchild[i] = n.children[i]
		}
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
	switch len(n.children) {
	case -1:
		// never 16
		idx := byte(n.size)
		bitfield := uint(0)
		for i := 0; i < int(n.size); i++ {
			if n.children[i].key[0] >= k {
				bitfield |= (1 << i)
			}
		}
		mask := (1 << len(n.children)) - 1
		bitfield &= uint(mask)
		if bitfield != 0 {
			idx = byte(bits.TrailingZeros(bitfield))
		}
		if idx < byte(n.size) && n.children[idx].key[0] == k {
			found = true
		}
		return int16(idx), found
	default:

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
	n.children = make([]*node, 4)
	// add childs
	n.children[0] = &node{
		key: n.key,
		val: n.val,
	}
	// modify to node with childs
	n.key = cp
	n.val = nil
	//ind, f := n.find(key[depth])
	//fmt.Println(ind, f)
	n.add(key[depth:], val, 0)
}

func (n *node) Print(depth int) {
	fmt.Print(fmt.Sprintf("%s %+v\n", strings.Repeat(" ", depth), n))
	depth++
	for idx := 0; idx < int(n.size); idx++ {
		n.children[idx].Print(depth)
	}
}

func (a *Art) Print() {
	a.root.Print(0)
}

/*
func (n *node) clone(depth int) (nn *node) {
	fmt.Print(fmt.Sprintf("%s %+v\n", strings.Repeat(" ", depth), n))
	depth++
	for idx := 0; idx < int(n.size); idx++ {
		n.children[idx].Print(depth)
	}
}
*/
