package files

import (
	"git-split/graphs"
	"strings"
)

type DependencySearchQuery struct {
	Node             *graphs.Node
	fileSearcheScore map[string]float64
}

func GetFileDependenciesQuery(filesGraph *FilesGraph) []DependencySearchQuery {
	allFiles := filesGraph.GetLeaves()
	var dependencySearches []DependencySearchQuery
	for _, file := range allFiles {
		fileSearches := getPossibleFileRefs(file)
		dependencySearches = append(dependencySearches, DependencySearchQuery{
			Node:             file,
			fileSearcheScore: fileSearches,
		})
	}
	return dependencySearches
}

func getPossibleFileRefs(node *graphs.Node) map[string]float64 {
	allRefs := make(map[string]float64)
	path := node.ID
	path = strings.TrimPrefix(path, "./")
	parts := strings.Split(path, "/")

	pathSplitOnExtension := strings.Split(path, ".")
	pathNoExtension := strings.Join(pathSplitOnExtension[:len(pathSplitOnExtension)-1], ".")
	partsNoExtension := strings.Split(pathNoExtension, "/")

	fileName := parts[len(parts)-1]
	baseName := fileName
	if idx := strings.LastIndex(fileName, "."); idx != -1 {
		baseName = fileName[:idx]
	}
	length := len(parts)
	for i := range parts {
		weight := calcWeight(length-i, length)
		addWithSeparators(allRefs, strings.Join(parts[i:], "/"), weight)
		addWithSeparators(allRefs, strings.Join(partsNoExtension[i:], "/"), weight)
	}
	addWithSeparators(allRefs, fileName, 0.6)
	addWithSeparators(allRefs, baseName, 0.5)
	return allRefs
}

func calcWeight(depth, length int) float64 {
	weight := 1.0
	switch depth {
	// if it is the full path or close to it, it is very likely the file is referenced, but could be comment / string / other file reference
	case 0:
		weight = 0.95
	case 1:
		weight = 0.92
	case 2:
		weight = 0.9

	// if it is only the file name or a couple of parent directories, it is likely the file is referenced, but could be comment / string / other file reference
	case length - 1:
		weight = 0.7
	case length - 2:
		weight = 0.8
	case length - 3:
		weight = 0.75

	// if it is somewhere inbetween, most likely is not a reference
	default:
		weight = 0.1
	}
	return weight
}

func addWithSeparators(refs map[string]float64, ref string, weight float64) {
	separators := []string{"/", ".", "\\", "\\\\", "//"}
	parts := strings.Split(ref, "/")
	for _, sep := range separators {
		joined := strings.Join(parts, sep)
		if existing, ok := refs[joined]; !ok || weight > existing {
			refs[joined] = weight
		}
	}
}
