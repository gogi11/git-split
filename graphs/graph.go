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
func (g *Graph) AddEdge(from, to, typ, label string, weight float64) {
	edge := &Edge{
		From:   from,
		To:     to,
		Type:   typ,
		Label:  label,
		Weight: weight,
		Attrs:  map[string]string{},
	}

	g.Edges = append(g.Edges, edge)

	g.Outgoing[from] = append(g.Outgoing[from], edge)
	g.Incoming[to] = append(g.Incoming[to], edge)
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

func (g *Graph) GetRoots() []*Node {
	var roots []*Node
	for _, node := range g.Nodes {
		if len(g.Incoming[node.ID]) == 0 {
			roots = append(roots, node)
		}
	}
	return roots
}
