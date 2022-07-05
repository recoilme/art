package art

import (
	"bytes"
)

type Art struct {
	root *node
}

func New() *Art {
	return &Art{}
}

func (a *Art) Set(key, val []byte) {
	//fmt.Println("Set", key)
	if a.root == nil {
		a.root = &node{
			key: key,
			val: val,
		}
		return
	}
	a.root.set(key, val, 0)
	return
}

func (a *Art) Get(key []byte) (val []byte) {
	//fmt.Println("Get", key)
	if a.root == nil {
		return nil
	}
	return a.root.get(key, 0)
}

func (a *Art) String() string {
	if a.root == nil {
		return ""
	}
	return a.root.String(0)
}

func (a *Art) StringKeys(isString bool) string {
	if a.root == nil {
		return ""
	}
	return "\n" + a.root.StringKeys(0, isString)
}

func (a *Art) Delete(key []byte) {
	if a.root == nil {
		return
	}
	// node without children
	if a.root.size == 0 && bytes.Equal(key, a.root.key) {
		a.root = nil
		return
	}

	for {
		key = a.root.delete(key, 0)
		if key == nil {
			break
		}
	}
	// set to nil if needed
	if a.root.val == nil && a.root.size == 0 {
		a.root = nil
	}

	return
}
