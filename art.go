package art

import "bytes"

type Art struct {
	root *node
}

type node struct {
	key []byte
	val []byte
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
	replaced = a.root.set(key, val)
	if replaced {
		return
	}
	return
}

func (n *node) set(key, val []byte) (replaced bool) {
	i, found := n.find(key)
	if found {
		if i < 0 {
			n.val = val
		}
		replaced = true
		return
	}
	return
}

func (n *node) find(key []byte) (index int, found bool) {
	if n.val != nil {
		if bytes.Equal(n.key, key) {
			index = -1
			found = true
		}
		return
	}
	return
}
