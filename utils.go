package art

import (
	"bytes"
	"fmt"
	"strings"
)

func commonPrefix(key, newKey []byte) []byte {
	i := 0
	limit := min(len(key), len(newKey))
	for ; i < limit; i++ {
		if key[i] != newKey[i] {
			if len(key[:i]) == 0 {
				return nil
			}
			return key[:i]
		}
	}
	if len(key[:i]) == 0 {
		return nil
	}
	return key[:i]
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func (n *node) String(depth int) string {
	buf := bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("%s%+v\n", strings.Repeat(" ", depth), n))
	depth++

	if len(n.children) == 256 {
		for idx := 0; idx < len(n.children); idx++ {
			if n.children[idx] != nil {
				buf.WriteString(n.children[idx].String(depth))
			}
		}
	} else {
		for idx := 0; idx < int(n.size); idx++ {
			buf.WriteString(n.children[idx].String(depth))
		}
	}
	return buf.String()
}

func (n *node) StringKeys(depth int, isString bool) string {
	buf := bytes.Buffer{}
	if isString {
		buf.WriteString(fmt.Sprintf("%skey:%s val:%s\n", strings.Repeat(" ", depth), n.key, n.val))
	} else {
		buf.WriteString(fmt.Sprintf("%skey:%v val:%v\n", strings.Repeat(" ", depth), n.key, n.val))
	}
	depth++
	for idx := 0; idx < int(n.size); idx++ {
		buf.WriteString(n.children[idx].StringKeys(depth, isString))
	}
	return buf.String()
}
