package graphviz

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"git-split/graphs"
)

func ToDOT(g *graphs.Graph) string {
	var b strings.Builder

	b.WriteString("digraph G {\n")
	b.WriteString("  rankdir=LR;\n")
	b.WriteString("  node [shape=box style=filled];\n\n")

	// Nodes
	for _, n := range g.Nodes {
		label := n.Label
		if label == "" {
			label = n.ID
		}
		attrs := []string{
			fmt.Sprintf(`label="%s"`, label),
		}
		for k, v := range n.Attrs {
			if k[0] == '_' {
				continue // skip internal attributes
			}
			if k == "change" {
				switch v {
				case "A":
					attrs = append(attrs, `fillcolor="palegreen"`)
				case "M":
					attrs = append(attrs, `fillcolor="khaki"`)
				case "D":
					attrs = append(attrs, `fillcolor="lightcoral"`)
				case "R":
					attrs = append(attrs, `fillcolor="lightblue"`)
				default:
					attrs = append(attrs, `fillcolor="white"`)
				}
				fmt.Printf("Node %s has change type %s\n", n.ID, v)
			} else {
				attrs = append(attrs, fmt.Sprintf(`%s="%s"`, k, v))
			}
		}
		attrString := strings.Join(attrs, ",")
		b.WriteString(fmt.Sprintf(`  "%s" [%s];`, n.ID, attrString))
		b.WriteString("\n")
	}
	b.WriteString("\n")

	// Edges
	for _, e := range g.Edges {
		attrs := []string{}
		if e.Label != "" {
			attrs = append(attrs, fmt.Sprintf(`label="%s"`, e.Label))
		}
		if e.Type != "" {
			attrs = append(attrs, fmt.Sprintf(`type="%s"`, e.Type))
		}
		if e.Weight != 0 {
			attrs = append(attrs, fmt.Sprintf(`weight="%f"`, e.Weight))
		}
		for k, v := range e.Attrs {
			attrs = append(attrs, fmt.Sprintf(`%s="%s"`, k, v))
		}
		attrString := ""
		if len(attrs) > 0 {
			attrString = " [" + strings.Join(attrs, ",") + "]"
		}
		b.WriteString(fmt.Sprintf(
			`  "%s" -> "%s"%s;`,
			e.From,
			e.To,
			attrString,
		))
		b.WriteString("\n")
	}
	b.WriteString("}\n")
	return b.String()
}

func runGraphviz(args ...string) (string, error) {
	cmd := exec.Command("dot", args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("graphviz %v failed: %v\n%s",
			args,
			err,
			stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func WriteDOTFile(g *graphs.Graph, path string, verbose bool) error {
	dot := ToDOT(g)
	if verbose {
		fmt.Printf("Writing DOT file to %s\n", path)
		fmt.Printf("Dot: \n%s", dot)
	}
	return os.WriteFile(path, []byte(dot), 0644)
}

func RenderDOTToFile(dotPath, outputPath, format string) error {
	_, err := runGraphviz("-T"+format, dotPath, "-o", outputPath)
	return err
}

func CreateGraphImage(g *graphs.Graph, verbose bool) {
	WriteDOTFile(g, "output/graph.dot", verbose)
	RenderDOTToFile("output/graph.dot", "output/graph.png", "png")
}
