package helpers

func Unique(input []string) []string {

	seen := map[string]bool{}
	var result []string

	for _, v := range input {

		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}

	return result
}
