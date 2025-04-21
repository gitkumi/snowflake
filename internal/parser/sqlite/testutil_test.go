package parser_test

func compareStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	aCopy := make([]string, len(a))
	bCopy := make([]string, len(b))

	copy(aCopy, a)
	copy(bCopy, b)

	for _, aVal := range aCopy {
		found := false
		for i, bVal := range bCopy {
			if aVal == bVal {
				bCopy = append(bCopy[:i], bCopy[i+1:]...)
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}
