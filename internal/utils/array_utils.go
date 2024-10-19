package utils

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return IndexOf(a, x) >= 0
}

func IndexOf(a []string, x string) int {
	for index, n := range a {
		if x == n {
			return index
		}
	}
	return -1
}

func Remove(a []string, i int) []string {
	a[i] = a[len(a)-1]
	return a[:len(a)-1]
}

func RetainAll(cleaned *[]string, source []string) {
	for _, obj := range *cleaned {
		if !Contains(source, obj) {
			i := IndexOf(*cleaned, obj)
			*cleaned = Remove(*cleaned, i)
		}
	}
}
