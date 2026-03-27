package clustering

import (
	"git-split/graphs"
)

func clusterByParent(g *graphs.Graph, maxFilesCount int) map[string][]*graphs.Node {
	fileNodes := g.GetLeaves()

	clusterMap := map[string][]*graphs.Node{}
	for _, f := range fileNodes {
		parent := g.GetParent(f, "contains")
		if parent == nil {
			parent = g.GetRoots()[0]
		}
		parentEdges := g.Outgoing[parent.ID]
		for parentEdges != nil && len(parentEdges) > maxFilesCount && parent.ID != "." {
			parent = g.GetParent(parent, "contains")
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
	// cluster 1 id -> cluster 2 id -> weight
	clusterDeps := []ClusterWeightRelationship{}
	for startKey, startCluster := range clusterMap {
		for endKey, endCluster := range clusterMap {

			skip := startKey == endKey // skip same or already processed clusters
			for _, dep := range clusterDeps {
				if dep.start == startKey && dep.end == endKey || dep.start == endKey && dep.end == startKey {
					skip = true
					break
				}
			}
			if skip {
				continue
			}

			clusterDeps = append(clusterDeps, ClusterWeightRelationship{start: startKey, end: endKey, weight: 0})
			lastIndex := len(clusterDeps) - 1
			for _, startNode := range startCluster {
				for _, endNode := range endCluster {
					edge := g.GetEdge(startNode, endNode, "depends_on")
					if edge != nil {
						clusterDeps[lastIndex].weight += edge.Weight
					}
					edge = g.GetEdge(endNode, startNode, "depends_on")
					if edge != nil {
						clusterDeps[lastIndex].weight += edge.Weight
					}
				}
			}
		}
	}
	return clusterDeps
}

func HybridCluster(g *graphs.Graph, depThreshold float64, maxFilesCount int) [][]*graphs.Node {
	// create initial clusters by parent directories
	clusterMap := clusterByParent(g, maxFilesCount)

	// collect the total dependency weight between clusters
	clusterDependencies := getTotalDependencyWeightBetweenClusters(g, clusterMap)

	// remove all clusterDependencies with weight < depThreshold
	for i := len(clusterDependencies) - 1; i >= 0; i-- {
		if clusterDependencies[i].weight < depThreshold {
			clusterDependencies = append(clusterDependencies[:i], clusterDependencies[i+1:]...)
		}
	}

	// merge all the clusters that are left with high dependencies inbetween
	newClusterMap := map[string][]*graphs.Node{}
	for _, dep := range clusterDependencies {
		if newClusterMap[dep.start] == nil {
			// copy the cluster to the new map
			newClusterMap[dep.start] = make([]*graphs.Node, 0)
			newClusterMap[dep.start] = append(newClusterMap[dep.start], clusterMap[dep.start]...)
		}
		newClusterMap[dep.start] = append(newClusterMap[dep.start], clusterMap[dep.end]...)
	}
	// add all the unchanged clusters
	for i, cluster := range clusterMap {
		if newClusterMap[i] == nil {
			newClusterMap[i] = cluster
		}
	}

	result := [][]*graphs.Node{}
	for _, cluster := range newClusterMap {
		result = append(result, cluster)
	}
	return result
}
