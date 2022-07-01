package art

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
