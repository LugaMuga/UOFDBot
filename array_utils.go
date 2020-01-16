package main

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return indexOf(a, x) >= 0
}

func indexOf(a []string, x string) int {
	for index, n := range a {
		if x == n {
			return index
		}
	}
	return -1
}

func remove(a []string, i int) []string {
	a[i] = a[len(a)-1]
	return a[:len(a)-1]
}

func retainAll(cleaned *[]string, source []string) {
	for _, obj := range *cleaned {
		if !contains(source, obj) {
			i := indexOf(*cleaned, obj)
			*cleaned = remove(*cleaned, i)
		}
	}
}
