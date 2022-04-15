package graph

type Topo struct {
	topoVertexCount int
	inDegrees       map[string]int
	topoGraph       *Graph
	graph           *Graph
}

func (t Topo) CircleVertexIndex() *Vertex {
	// if len(t.topoGraph.graph) < t.topoVertexCount {
	for i, inDegree := range t.inDegrees {
		if inDegree > 0 {
			return t.graph.GetVertex(i)
		}
	}
	//}
	return nil
}
