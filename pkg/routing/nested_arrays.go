package routing

// GetDifferenceOfNestedArrays returns the difference between two deduplicated arrays
func GetDifferenceOfNestedArrays(old, new [][2]string) (deletions [][2]string, additions [][2]string) {
	for _, el := range old {
		exists := false
		for _, nel := range new {
			if (nel[0] == el[1] && nel[1] == el[0]) || (nel[0] == el[0] && nel[1] == el[1]) {
				exists = true

				break
			}
		}

		if !exists {
			deletions = append(deletions, el)
		}
	}

	for _, nel := range new {
		exists := false
		for _, el := range old {
			if (el[0] == nel[1] && el[1] == nel[0]) || (el[0] == nel[0] && el[1] == nel[1]) {
				exists = true

				break
			}
		}

		if !exists {
			additions = append(additions, nel)
		}
	}

	return deletions, additions
}

// GetUniqueKeys return the unique keys of a nested array
func GetUniqueKeys(in [][2]string) []string {
	outMap := make(map[string]bool)
	for _, key := range in {
		outMap[key[0]] = true
		outMap[key[1]] = true
	}

	var out []string
	for key := range outMap {
		out = append(out, key)
	}

	return out
}

// DeduplicateNestedArray deduplicates a nested array
func DeduplicateNestedArray(in [][2]string) [][2]string {
	var out [][2]string

	for _, el := range in {
		match := false
		for _, nel := range out {
			if (nel[0] == el[1] && nel[1] == el[0]) || (nel[0] == el[0] && nel[1] == el[1]) {
				match = true

				break
			}
		}

		if !match {
			out = append(out, el)
		}
	}

	return out
}
