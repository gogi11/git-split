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

func (fileGraph *FilesGraph) AddDependencyEdges() {
	query := GetFileDependenciesQuery(fileGraph)
	for _, node := range fileGraph.GetLeaves() {
		for _, search := range query {
			if node.ID == search.Node.ID {
				continue
			}
			// open the file and search for references to the changed file, if found add a dependency edge with weight based on how closely it matches the file path
			// search only the first 300 lines of the file to avoid performance issues, as references are likely to be found in imports / requires / file references at the top of the file
			for ref, weight := range search.fileSearcheScore {
				if strings.Contains(node.ID, ref) {
					fileGraph.AddEdgeOrMaxWeight(node.ID, search.Node.ID, "depends_on", "Dependency Score", weight)
				}
			}
		}
	}

	// all files (leaves) with same parent should have some depency weight between them, as they are likely to be related, even if no direct reference is found, so add a small weight to all leaves with same parent
	for _, node := range fileGraph.GetLeaves() {
		for _, sibling := range fileGraph.GetLeaves() {
			if node.ID != sibling.ID && strings.HasPrefix(node.ID, strings.TrimSuffix(sibling.ID, sibling.ID[strings.LastIndex(sibling.ID, "/"):])) {
				fileGraph.AddEdgeOrAccumulateWeight(node.ID, sibling.ID, "depends_on", "Sibling Dependency", 0.05)
			}
		}
	}
}
