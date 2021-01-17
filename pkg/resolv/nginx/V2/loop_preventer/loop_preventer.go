package loop_preventer

import "github.com/ClessLi/bifrost/pkg/graph"

type LoopPreventer interface {
	AddStringElement(src, dst string) error
}

type loopPreventer struct {
	graph *graph.Graph
}

func (l *loopPreventer) AddStringElement(src, dst string) error {
	return l.graph.AddEdge(src, dst)
}

func NewLoopPreverter(strElement string) LoopPreventer {
	return &loopPreventer{graph: graph.NewGraph(strElement)}
}
