package art

import (
	"bytes"
	"math/bits"
)

type node struct {
	key      []byte
	val      []byte
	children []*node
	size     int16
}

func (n *node) set(key, val []byte, depth int) (replaced bool) {
	// leaf node
	if n.size == 0 {
		if bytes.Equal(key[depth:], n.key) {
			n.val = val
			return true
		}
		n.nodeSplit(key[depth:], val)
		return false
	}
	// nodes with children
	if bytes.Equal(key[depth:], n.key) {
		// prefix equal
		n.val = val
		return true
	}
	insert := func() (inserted bool) {
		i := n.find(key[depth])
		if i >= 0 {
			n = n.children[i]
			n.set(key, val, depth)
			// go to child
		} else {
			n.add(key[depth:], val)
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

func (n *node) get(key []byte, depth int) (val []byte) {
	//fmt.Println("get", key, depth, val)
	// leaf node
	if n.size == 0 {
		if bytes.Equal(key[depth:], n.key) {
			val = n.val
			return
		}
		return nil
	}
	// nodes with children
	if bytes.Equal(key[depth:], n.key) {
		// prefix equal
		val = n.val
		return
	}
	// node without prefix
	if len(n.key) == 0 {
		i := n.find(key[depth])
		if i >= 0 {
			n = n.children[i]
			if bytes.Equal(key[depth:], n.key) {
				val = n.val
				return
			}
			return n.get(key, depth)
			// go to child
		} else {
			return nil
		}
	}
	// node with prefix
	cp := commonPrefix(key[depth:], n.key)
	delta := len(n.key) - len(cp)
	if delta <= 0 {
		depth += len(cp)
		i := n.find(key[depth])
		if i >= 0 {
			n = n.children[i]
			if bytes.Equal(key[depth:], n.key) {
				//fmt.Println("eq", key, n.val)
				val = n.val
				return
			}
			return n.get(key, depth)
			// go to child
		} else {
			return nil
		}
	}
	return
}

func (n *node) grow() {
	switch n.size {
	case 4, 16, 48:
		var newsize int
		switch n.size {
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

func (n *node) add(key, val []byte) {
	//	fmt.Println("add", key, n.size)
	n.grow()
	var idx int16
	for ; idx < n.size; idx++ {
		if key[0] < n.children[idx].key[0] {
			break
		}
	}

	for i := n.size; i > idx; i-- {
		if n.children[i-1].key[0] > key[0] {
			n.children[i] = n.children[i-1]
		}
	}
	n.children[idx] = &node{
		key: key,
		val: val,
	}
	n.size += 1
}

func (n *node) find(k byte) (index int16) {
	// TODO 255
	switch len(n.children) {
	case -1:
		// never 16
		bitfield := uint(0)
		for i := 0; i < int(n.size); i++ {
			if n.children[i].key[0] >= k {
				bitfield |= (1 << i)
			}
		}
		mask := (1 << len(n.children)) - 1
		bitfield &= uint(mask)
		if bitfield != 0 {
			return int16(byte(bits.TrailingZeros(bitfield)))
		}
		return -1
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
			return low - 1
		}
		return -1
	}
}

func (n *node) nodeSplit(key, val []byte) {
	//fmt.Println("nodeSplit", key)
	cp := commonPrefix(key, n.key)
	depth := len(cp)
	// add childs
	n.children = make([]*node, 4)
	if bytes.Equal(n.key, cp) {
		// old key = cp
		n.add(key[depth:], val)
		return
	}
	if bytes.Equal(key, cp) {
		// new key = cp
		oldkey := n.key
		oldval := n.val
		n.key = cp
		n.val = val
		n.add(oldkey[depth:], oldval)
		return
	}
	// save old vals
	n.children[0] = &node{
		key: n.key[depth:],
		val: n.val,
	}
	n.size = 1
	// modify
	n.key = cp
	n.val = nil
	n.add(key[depth:], val)
}
