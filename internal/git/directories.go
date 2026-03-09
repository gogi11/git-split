package git

import (
	"path/filepath"
)

// GroupFilesByDepthMap groups files by directory up to `depth` segments.
// Root-level files (no directory) are grouped under ".".
// Returns a map: groupDir -> list of files.
func GroupFilesByDepthMap(files []FileChange, depth int) map[string][]FileChange {
	result := make(map[string][]FileChange)
	for _, f := range files {
		path := filepath.Clean(f.Path)
		segments := splitPathSegments(path)
		var group string
		if len(segments) == 1 {
			group = "."
		} else {
			segDepth := depth
			if len(segments) < depth {
				segDepth = len(segments) - 1
			}
			group = filepath.Join(segments[:segDepth]...)
		}
		result[group] = append(result[group], f)
	}

	return result
}

func splitPathSegments(p string) []string {
	var segments []string
	p = filepath.Clean(p)
	for {
		dir, file := filepath.Split(p)
		if file != "" {
			segments = append([]string{file}, segments...)
		}
		if dir == "" || dir == "/" || dir == "." {
			break
		}
		p = dir[:len(dir)-1] // remove trailing separator
	}
	return segments
}
