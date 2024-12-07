package tool

// CheckSliceEqual 检查两个slice是否元素相同
func CheckSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	mp := make(map[string]int, len(a)+len(b))
	for _, v := range a {
		mp[v]++
	}
	for _, v := range b {
		mp[v]--
		if mp[v] < 0 {
			return false
		}
	}
	return true
}
