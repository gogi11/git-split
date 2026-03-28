package clustering

import (
	"git-split/graphs"
)

func clusterByParent(g *graphs.Graph, maxFilesCount int) map[string][]*graphs.Node {
	fileNodes := g.GetNodesByType("file")
	clusterMap := map[string][]*graphs.Node{}
	for _, f := range fileNodes {
		parent := g.GetParent(f, "contains")
		if parent == nil {
			roots := g.GetRoots()
			if len(roots) == 0 {
				continue
			}
			parent = roots[0]
		}
		// climb up until directory is "big enough"
		for parent != nil && parent.ID != "." {
			children := g.GetSubNodes(parent, "file")
			// stop if this dir is large enough
			if len(children) >= maxFilesCount {
				break
			}
			next := g.GetParent(parent, "contains")
			if next == nil {
				break
			}
			parent = next
		}
		clusterMap[parent.ID] = append(clusterMap[parent.ID], f)
	}
	return clusterMap
}

type ClusterWeightRelationship struct {
	start  string
	end    string
	weight float64
}

func getTotalDependencyWeightBetweenClusters(g *graphs.Graph, clusterMap map[string][]*graphs.Node) []ClusterWeightRelationship {
	clusterDeps := []ClusterWeightRelationship{}
	seen := map[string]bool{}
	for startKey, startCluster := range clusterMap {
		for endKey, endCluster := range clusterMap {
			if startKey == endKey {
				continue
			}

			// always start from the smaller name to the bigger (parent -> child, as these are paths)
			a, b := startKey, endKey
			if a > b {
				a, b = b, a
			}

			key := a + "::" + b
			if seen[key] {
				continue
			}
			seen[key] = true
			totalWeight := 0.0
			for _, startNode := range startCluster {
				for _, endNode := range endCluster {
					if edge := g.GetEdge(startNode, endNode, "depends_on"); edge != nil {
						totalWeight += edge.Weight
					}
					if edge := g.GetEdge(endNode, startNode, "depends_on"); edge != nil {
						totalWeight += edge.Weight
					}
				}
			}

			if totalWeight > 0 {
				clusterDeps = append(clusterDeps, ClusterWeightRelationship{
					start:  a,
					end:    b,
					weight: totalWeight,
				})
			}
		}
	}
	return clusterDeps
}

func cleanDependencies(clusterDependencies []ClusterWeightRelationship, depThreshold float64) []ClusterWeightRelationship {
	result := []ClusterWeightRelationship{}
	for _, dep := range clusterDependencies {
		if dep.weight >= depThreshold {
			result = append(result, dep)
		}
	}
	return result
}

func mergeClusters(
	clusterDependencies []ClusterWeightRelationship, allClusters map[string][]*graphs.Node) map[string][]*graphs.Node {
	// Do a DFS to find connected clusters and merge them together

	// 1. Make adjacency list for connected clusters (undirected)
	adj := make(map[string][]string)
	for _, dep := range clusterDependencies {
		adj[dep.start] = append(adj[dep.start], dep.end)
		adj[dep.end] = append(adj[dep.end], dep.start)
	}

	visited := make(map[string]bool)
	result := make(map[string][]*graphs.Node)
	// 2. Traverse clusters (connected components)
	for clusterID := range allClusters {
		if visited[clusterID] {
			continue
		}
		// BFS/DFS to collect component
		stack := []string{clusterID}
		component := []string{}

		for len(stack) > 0 {
			curr := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if visited[curr] {
				continue
			}
			visited[curr] = true
			component = append(component, curr)
			for _, nei := range adj[curr] {
				if !visited[nei] {
					stack = append(stack, nei)
				}
			}
		}

		// 3. Merge all clusters in this component
		root := component[0]
		merged := []*graphs.Node{}
		seenNodes := map[string]bool{}
		for _, cid := range component {
			for _, n := range allClusters[cid] {
				if !seenNodes[n.ID] {
					merged = append(merged, n)
					seenNodes[n.ID] = true
				}
			}
		}
		result[root] = merged
	}

	return result
}

func HybridCluster(g *graphs.Graph, depThreshold float64, maxFilesCount int) [][]*graphs.Node {
	// create initial clusters by parent directories
	clusterMap := clusterByParent(g, maxFilesCount)

	// collect the total dependency weight between clusters
	clusterDependencies := getTotalDependencyWeightBetweenClusters(g, clusterMap)

	// remove all clusterDependencies with weight < depThreshold
	clusterDependencies = cleanDependencies(clusterDependencies, depThreshold)

	// merge clusters with inbetween dependencies
	newClusterMap := mergeClusters(clusterDependencies, clusterMap)

	result := [][]*graphs.Node{}
	for _, cluster := range newClusterMap {
		result = append(result, cluster)
	}
	return result
}
