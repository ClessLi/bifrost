package graph

import "github.com/ClessLi/component-base/pkg/errors"

var (
	ErrVertexRepeat           = errors.New("vertex repeat error")
	ErrVertexNotExist         = errors.New("vertex not exist")
	ErrEdgeRepeat             = errors.New("edge repeat error")
	ErrLoopOfTopologicalGraph = errors.New("loop of topological graph")
)
