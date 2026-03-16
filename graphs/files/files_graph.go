package files

import (
	"fmt"
	"strings"

	"git-split/graphs"
)

type FilesGraph struct {
	*graphs.Graph
}

func NewFilesGraph(actions []string, paths [][]string) *FilesGraph {
	fileGraph := &FilesGraph{graphs.NewGraph()}
	fileGraph.AddNode(".", ".", "directory")
	fileGraph.Nodes["."].Attrs["depth"] = "0"
	for i, path := range paths {
		var oldPath string
		for j, p := range path {
			dirs := strings.Split(p, "/")
			currentDir := "."
			for depth, dirName := range dirs {
				parentDir := currentDir
				currentDir = strings.TrimRight(currentDir, "/") + "/" + dirName
				if currentDir != "./"+p { // if directory, add the node
					fileGraph.AddNode(currentDir, dirName, "directory")
				} else { // if file, add the node and edge based on action
					fileGraph.AddNode(currentDir, dirName, "file")
					if actions[i] != "R" { // add edge from parent to file if not a rename
						fileGraph.AddEdge(parentDir, currentDir, actions[i], "", 1)
					} else if j == 1 { // if it is a rename (move) and is new name add edge from old path to new path
						fileGraph.AddEdge(oldPath, currentDir, "R", "", 1)
					}
					fileGraph.Nodes[currentDir].Attrs["depth"] = fmt.Sprintf("%d", depth)
					fileGraph.Nodes[currentDir].Attrs["change"] = actions[i]
				}
				fileGraph.AddEdge(parentDir, currentDir, "contains", "", 1)
			}
			oldPath = currentDir
		}
	}
	return fileGraph
}
