package routing

// GetDifferenceOfNestedArrays returns the difference between two deduplicated arrays
func GetDifferenceOfNestedArrays(old, new [][]string) (deletions [][]string, additions [][]string) {
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
func GetUniqueKeys(in [][]string) []string {
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
