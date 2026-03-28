package graphs

type Node struct {
	ID    string
	Label string
	Type  string

	Attrs map[string]string
}

type Edge struct {
	From string
	To   string

	Type  string
	Label string

	Weight float64

	Attrs map[string]string
}

type Graph struct {
	Nodes map[string]*Node
	Edges []*Edge

	Outgoing map[string][]*Edge
	Incoming map[string][]*Edge
}

func NewGraph() *Graph {
	return &Graph{
		Nodes:    make(map[string]*Node),
		Edges:    []*Edge{},
		Outgoing: make(map[string][]*Edge),
		Incoming: make(map[string][]*Edge),
	}
}

func (g *Graph) AddNode(id, label, typ string) {
	g.Nodes[id] = &Node{
		ID:    id,
		Label: label,
		Type:  typ,
		Attrs: map[string]string{},
	}
}
func (g *Graph) addEdge(from, to, typ, label string, weight float64, handleDuplicate string) {
	edge := &Edge{
		From:   from,
		To:     to,
		Type:   typ,
		Label:  label,
		Weight: weight,
		Attrs:  map[string]string{},
	}

	if g.Outgoing[from] != nil {
		for _, e := range g.Outgoing[from] {
			if e.To == to && e.Type == typ {
				switch handleDuplicate {
				case "accumulate":
					e.Weight += weight
				case "max":
					if e.Weight < weight {
						e.Weight = weight
					}
				default:
					return
				}
				return
			}
		}
	}
	g.Edges = append(g.Edges, edge)

	g.Outgoing[from] = append(g.Outgoing[from], edge)
	g.Incoming[to] = append(g.Incoming[to], edge)
}

func (g *Graph) AddEdge(from, to, typ, label string, weight float64) {
	g.addEdge(from, to, typ, label, weight, "")
}
func (g *Graph) AddEdgeOrAccumulateWeight(from, to, typ, label string, weight float64) {
	g.addEdge(from, to, typ, label, weight, "accumulate")
}
func (g *Graph) AddEdgeOrMaxWeight(from, to, typ, label string, weight float64) {
	g.addEdge(from, to, typ, label, weight, "max")
}

func (g *Graph) RemoveEdge(edge *Edge) {
	for i, e := range g.Edges {
		if e == edge {
			g.Edges = append(g.Edges[:i], g.Edges[i+1:]...)
			break
		}
	}
	for i, e := range g.Outgoing[edge.From] {
		if e == edge {
			g.Outgoing[edge.From] = append(g.Outgoing[edge.From][:i], g.Outgoing[edge.From][i+1:]...)
			if len(g.Outgoing[edge.From]) == 0 {
				delete(g.Outgoing, edge.From)
			}
			break
		}
	}
	for i, e := range g.Incoming[edge.To] {
		if e == edge {
			g.Incoming[edge.To] = append(g.Incoming[edge.To][:i], g.Incoming[edge.To][i+1:]...)
			if len(g.Incoming[edge.To]) == 0 {
				delete(g.Incoming, edge.To)
			}
			break
		}
	}
}

func (g *Graph) GetParent(node *Node, edgeType string) *Node {
	for _, edge := range g.Incoming[node.ID] {
		if edge.Type == edgeType {
			return g.Nodes[edge.From]
		}
	}
	return nil
}

func (g *Graph) GetEdge(start, end *Node, edgeType string) *Edge {
	for _, edge := range g.Outgoing[start.ID] {
		if edge.Type == edgeType {
			if edge.To == end.ID {
				return edge
			}
		}
	}
	return nil
}

func (g *Graph) FilterEdges(typ string, weight float64) []*Edge {
	// Removes all edges of the type typ and weight less than weight
	// and returns all removed edges
	var edges []*Edge
	for _, edge := range g.Edges {
		if edge.Type == typ && edge.Weight < weight {
			edges = append(edges, edge)
		}
	}
	for _, edge := range edges {
		g.RemoveEdge(edge)
	}
	return edges
}

func (g *Graph) GetNode(id string) (*Node, bool) {
	node, exists := g.Nodes[id]
	return node, exists
}

func (g *Graph) GetLeaves() []*Node {
	var leaves []*Node
	for _, node := range g.Nodes {
		if len(g.Outgoing[node.ID]) == 0 {
			leaves = append(leaves, node)
		}
	}
	return leaves
}

func (g *Graph) GetSubNodes(parent *Node, nodeTyp string) []*Node {
	// do a DFS to find all nodes of type nodeTyp that are children of parent
	var nodes []*Node
	queue := []*Node{parent}
	seenNodes := map[string]bool{}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		if seenNodes[node.ID] {
			continue
		}
		seenNodes[node.ID] = true
		if node.Type == nodeTyp {
			nodes = append(nodes, node)
		}
		for _, edge := range g.Outgoing[node.ID] {
			queue = append(queue, g.Nodes[edge.To])
		}
	}
	return nodes
}

func (g *Graph) GetNodesByType(typ string) []*Node {
	var leaves []*Node
	for _, node := range g.Nodes {
		if node.Type == typ {
			leaves = append(leaves, node)
		}
	}
	return leaves
}

func (g *Graph) GetRoots() []*Node {
	var roots []*Node
	for _, node := range g.Nodes {
		if len(g.Incoming[node.ID]) == 0 {
			roots = append(roots, node)
		}
	}
	return roots
}
