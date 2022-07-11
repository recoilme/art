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

func (n *node) get(key []byte, depth int, strict bool) (*node, int) {
	//fmt.Println("get", key, depth, val)
	// node without children
	if n.size == 0 {
		if compare(key[depth:], n.key, strict) {
			return n, depth + len(commonPrefix(key[depth:], n.key))
		}
		return nil, 0
	}
	// nodes with children
	if compare(key[depth:], n.key, strict) {
		// prefix equal
		return n, depth + len(commonPrefix(key[depth:], n.key))
	}
	// node with or without prefix
	cp := commonPrefix(key[depth:], n.key)
	delta := len(n.key) - len(cp)
	if delta <= 0 {
		depth += len(cp)
		i := n.find(key[depth])
		if i >= 0 {
			n = n.children[i]
			if compare(key[depth:], n.key, strict) {
				return n, depth + len(commonPrefix(key[depth:], n.key)) //len(n.key)
			}
			// go to child
			return n.get(key, depth, strict)
		}
		return nil, 0
	}
	return nil, 0
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

func (n *node) scan(iter func(key, val []byte) bool, prefix []byte, depth int) bool {
	//fmt.Println("scan", depth, string(prefix), string(prefix[:depth]), string(n.key))
	if len(n.key) > 0 {
		prefix = append(prefix, n.key...)
		depth += len(n.key)
	}
	if n.val != nil {
		if !iter(prefix, n.val) {
			return false
		}
	}

	var i int16
	size := n.size
	if len(n.children) == 256 {
		size = 256
	}
	for ; i < size; i++ {
		if size == 256 && n.children[i] == nil {
			continue
		}
		if n.children[i].size == 0 {
			if len(prefix) > 0 {
				if !iter(append(prefix, n.children[i].key...), n.children[i].val) {
					break
				}
				continue
			}
			if !iter(n.children[i].key, n.children[i].val) {
				break
			}
			continue
		}
		n.children[i].scan(iter, prefix, depth)
	}
	if n.key != nil {
		depth -= len(n.key)
		if depth < 0 {
			prefix = []byte{}
			return false
		} else {
			prefix = prefix[:depth]
		}
	}
	return false
}

func (n *node) ascend(pivot []byte, iter func(key, val []byte) bool) bool {
	n, depth := n.get(pivot, 0, false)
	if n != nil {
		cs := commonSuffix(pivot[:depth], n.key)
		depth -= len(cs)
		pref := pivot[:depth]
		return n.scan(iter, pref, len(pref))
	}
	return false
}

func (n *node) descend(pivot []byte, iter func(key, val []byte) bool) bool {
	n, depth := n.get(pivot, 0, false)
	if n != nil {
		cs := commonSuffix(pivot[:depth], n.key)
		depth -= len(cs)
		pref := pivot[:depth]

		keys := make([][]byte, 0, 8)
		vals := make([][]byte, 0, 8)
		n.scan(func(key, val []byte) bool {
			newkey := make([]byte, len(key))
			copy(newkey, key)
			keys = append(keys, newkey)
			vals = append(vals, val)
			return true
		}, pref, len(pref))
		for i := len(keys) - 1; i >= 0; i-- {
			if !iter(keys[i], vals[i]) {
				break
			}
		}
	}
	return false
}

func (n *node) scanNodes(iter func(key []byte, nn *node) bool, prefix []byte, depth int) bool {
	if n.val != nil {
		if n.size == 0 {
			if !iter(append(prefix, n.key...), n) {
				return false
			}
		}
		if !iter(prefix, n) {
			return false
		}
	}

	var i int16
	size := n.size
	if len(n.children) == 256 {
		size = 256
	}
	for ; i < size; i++ {
		if size == 256 && n.children[i] == nil {
			continue
		}
		if n.children[i].size == 0 {
			if len(prefix) > 0 {
				if !iter(append(prefix, n.children[i].key...), n.children[i]) {
					break
				}
				continue
			}
			if !iter(n.children[i].key, n.children[i]) {
				break
			}
			continue
		}
		if len(n.children[i].key) > 0 {
			prefix = append(prefix, n.children[i].key...)
			depth += len(n.children[i].key)
		}
		n.children[i].scanNodes(iter, prefix, depth)
		depth -= len(n.children[i].key)
		prefix = prefix[:depth]
	}
	return false
	/*
		if n.val != nil {
			if !iter([]byte(prefix+string(n.key)), n) {
				return false
			}
		}
		prefix += string(n.key)
		if len(n.children) == 256 {
			for i := 0; i < len(n.children); i++ {
				if n.children[i] != nil {
					n.children[i].scanNodes(iter, prefix)
				}
			}
		} else {
			for i := 0; i < int(n.size); i++ {
				n.children[i].scanNodes(iter, prefix)
			}
		}
		return false
	*/
}
