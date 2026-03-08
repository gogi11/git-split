package git

import "strings"

func GroupFilesByDepth(files []string, depth int) []string {
	set := map[string]bool{}
	for _, f := range files {
		parts := strings.Split(f, "/")
		if len(parts) < depth {
			depth = len(parts)
		}
		group := strings.Join(parts[:depth], "/")
		set[group] = true
	}
	var result []string
	for k := range set {
		result = append(result, k)
	}
	return result
}
