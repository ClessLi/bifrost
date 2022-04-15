package loop_preventer

import (
	"fmt"

	"github.com/marmotedu/errors"

	"github.com/ClessLi/bifrost/pkg/graph"
)

type LoopPreventer interface {
	CheckLoopPrevent(src, dst string) error
	RemoveRoute(src, dst string) error
}

type loopPreventer struct {
	graph *graph.Graph
}

func (l *loopPreventer) CheckLoopPrevent(src, dst string) error {
	return l.graph.AddEdge(src, dst)
}

func (l *loopPreventer) RemoveRoute(src, dst string) error {
	err := l.graph.DelEdge(src, dst)
	if errors.Is(err, graph.ErrVertexNotExist) {
		return fmt.Errorf("position of %s is not exist", src)
	}

	return err
}

func NewLoopPreverter(strElement string) LoopPreventer {
	return &loopPreventer{graph: graph.NewGraph(strElement)}
}
