package art

import (
	"bytes"
)

type node struct {
	key      []byte
	val      []byte
	children []*node
	size     int16
}

func (n *node) set(key, val []byte, depth int) {
	// node without children
	if n.size == 0 {
		if bytes.Equal(key[depth:], n.key) {
			n.val = val
			return
		}
		n.nodeSplit(key[depth:], val)
		return
	}
	// nodes with children
	if bytes.Equal(key[depth:], n.key) {
		// prefix equal
		n.val = val
		return
	}
	insert := func() {
		i := n.find(key[depth])
		if i >= 0 {
			n = n.children[i]
			n.set(key, val, depth)
			// go to child
			return
		}
		n.add(key[depth:], val)
		return
	}

	cp := commonPrefix(key[depth:], n.key)
	delta := len(n.key) - len(cp)
	if delta <= 0 {
		depth += len(cp)
		insert()
	} else {
		// rebase tree
		oldnode := &node{}
		*oldnode = *n
		oldnode.key = oldnode.key[len(cp):]
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
		return
	}
	return
}

func (n *node) get(key []byte, depth int) []byte {
	//fmt.Println("get", key, depth, val)
	// node without children
	if n.size == 0 {
		if bytes.Equal(key[depth:], n.key) {
			return n.val
		}
		return nil
	}
	// nodes with children
	if bytes.Equal(key[depth:], n.key) {
		// prefix equal
		return n.val
	}
	// node with or without prefix
	cp := commonPrefix(key[depth:], n.key)
	delta := len(n.key) - len(cp)
	if delta <= 0 {
		depth += len(cp)
		i := n.find(key[depth])
		if i >= 0 {
			n = n.children[i]
			if bytes.Equal(key[depth:], n.key) {
				return n.val
			}
			// go to child
			return n.get(key, depth)
		}
		return nil
	}
	return nil
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
		if newsize == 256 {
			for i = 0; i < n.size; i++ {
				idx := n.children[i].key[0]
				newchild[idx] = n.children[i]
			}
		} else {
			for i = 0; i < n.size; i++ {
				newchild[i] = n.children[i]
			}
		}
		n.children = newchild
	}
}

func (n *node) add(key, val []byte) {
	//	fmt.Println("add", key, n.size)
	n.grow()
	if len(n.children) == 256 {
		n.children[key[0]] = &node{
			key: key,
			val: val,
		}
		n.size += 1
		return
	}
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
	switch len(n.children) {
	case 4, 16:
		var idx int16
		for ; idx < n.size; idx++ {
			if k == n.children[idx].key[0] {
				return idx
			}
		}
		return -1
	case 256:
		if n.children[k] != nil {
			return int16(k)
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
		n.add(key[depth:], val)
		return
	}
	if bytes.Equal(key, cp) {
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

func (n *node) delete(key []byte, depth int) []byte {
	//fmt.Println("delete", key)
	// nodes with children
	if bytes.Equal(key[depth:], n.key) {
		// prefix equal
		n.val = nil
		return nil
	}
	cp := commonPrefix(key[depth:], n.key)
	delta := len(n.key) - len(cp)
	if delta <= 0 {
		depth += len(cp)
		i := n.find(key[depth])
		if i >= 0 {
			child := n.children[i]
			if bytes.Equal(key[depth:], child.key) {
				n.del(i)
				if n.val == nil && n.size == 0 {
					// self remove
					return key[:depth]
				}
				return nil
			}
			// go to child
			return child.delete(key, depth)
		}
		return nil
	}
	return nil
}

func (n *node) del(idx int16) {
	//fmt.Println("del", idx)
	n.children[idx] = nil
	if len(n.children) != 256 {
		for i := idx; i < n.size-1; i++ {
			n.children[i] = n.children[i+1]
		}
	}

	n.size--
	if n.size == 0 {
		n.children = nil
	}
}

func (n *node) scan(iter func(key, val []byte) bool, prefix string) bool {

	if n != nil && n.val != nil && n.key != nil {
		if !iter([]byte(prefix+string(n.key)), n.val) {
			return false
		}
	}
	prefix += string(n.key)
	if len(n.children) == 256 {
		for i := 0; i < len(n.children); i++ {
			if n.children[i] != nil {
				n.children[i].scan(iter, prefix)
			}
		}
	} else {
		for i := 0; i < int(n.size); i++ {
			n.children[i].scan(iter, prefix)
		}
	}
	return false
}
