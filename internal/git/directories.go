package git

import "strings"

func GroupFilesByDirectory(files []string) map[string][]string {
	result := map[string][]string{}
	for _, f := range files {
		parts := strings.Split(f, "/")
		dir := parts[0]
		result[dir] = append(result[dir], f)
	}
	return result
}
