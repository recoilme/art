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

func (n *node) String(depth int) string {
	buf := bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("%s%+v\n", strings.Repeat(" ", depth), n))
	depth++
	for idx := 0; idx < int(n.size); idx++ {
		buf.WriteString(n.children[idx].String(depth))
	}
	return buf.String()
}

func (n *node) StringKeys(depth int, isString bool) string {
	buf := bytes.Buffer{}
	if isString {
		buf.WriteString(fmt.Sprintf("%s%s:%s\n", strings.Repeat(" ", depth), n.key, n.val))
	} else {
		buf.WriteString(fmt.Sprintf("%s%v:%v\n", strings.Repeat(" ", depth), n.key, n.val))
	}
	depth++
	for idx := 0; idx < int(n.size); idx++ {
		buf.WriteString(n.children[idx].StringKeys(depth, isString))
	}
	return buf.String()
}
