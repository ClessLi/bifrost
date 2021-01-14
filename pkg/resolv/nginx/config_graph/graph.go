package config_graph

import (
	"github.com/ClessLi/bifrost/pkg/queue"
)

type Graph struct {
	graph          map[string]*Vertex
	StartingVertex *Vertex
	currentVertex  *Vertex
}

func (g *Graph) SetCurrentVertex(id string) error {
	v := g.GetVertex(id)
	if v == nil {
		return ErrVertexNotExist
	}
	g.currentVertex = v
	return nil
}

func (g Graph) GetCurrentVertexId() string {
	return g.currentVertex.id
}

func (g *Graph) AddEdgeFromCurrentVertex(tid string) error {
	return g.AddEdge(g.currentVertex.id, tid)
}

func (g *Graph) AddEdge(sid, tid string) error {
	sv := g.GetVertex(sid)
	if sv == nil {
		g.graph[sid] = newVertex(sid)
		sv = g.GetVertex(sid)
	}
	err := sv.AddEdge(tid)
	if err != nil {
		return err
	}
	tv := g.GetVertex(tid)
	if tv == nil {
		g.graph[tid] = newVertex(tid)
	}
	return nil
}

func (g Graph) GetVertex(sid string) *Vertex {
	v, ok := g.graph[sid]
	if ok {
		return v
	}
	return nil
}

func (g *Graph) TopoLogicalSortByKahn() error {
	vertexCount := len(g.graph)
	inDegrees := make(map[string]int)
	for idx := range g.graph {
		inDegrees[idx] = 0
	}
	for _, vertex := range g.graph {
		for j := 0; j < int(vertex.Size()); j++ {
			tid := vertex.GetEdgeByIndex(j).tid
			if tid == "" {
				return nil
			}
			if _, ok := inDegrees[tid]; ok {
				inDegrees[tid]++
			} else {
				inDegrees[tid] = 1
			}
		}
	}

	stringQueue := make(queue.StringQueue, 0)
	for idx, inDegree := range inDegrees {
		if inDegree == 0 {
			stringQueue.Add(idx)
		}
	}
	if stringQueue.IsEmpty() {
		return ErrLoopOfTopologicalGraph
	}
	//stringQueue.Add(g.StartingVertex.id)

	var currentVertex *Vertex
	var topoGraph *Graph
	for !stringQueue.IsEmpty() {
		id := stringQueue.Poll()
		if currentVertex == nil || topoGraph == nil {
			topoGraph = NewGraph(id)
			currentVertex = topoGraph.StartingVertex
		} else {
			err := topoGraph.AddEdge(currentVertex.id, id)
			if err != nil {
				return nil
			}
			currentVertex = topoGraph.GetVertex(id)
		}
		for j := 0; j < int(g.graph[id].Size()); j++ {
			k := g.graph[id].GetEdgeByIndex(j).tid
			inDegrees[k]--
			if inDegrees[k] == 0 {
				stringQueue.Add(k)
			}
		}
	}
	topo := &Topo{
		topoVertexCount: vertexCount,
		inDegrees:       inDegrees,
		topoGraph:       topoGraph,
		graph:           g,
	}
	loopVertex := topo.CircleVertexIndex()
	if loopVertex != nil {
		return ErrLoopOfTopologicalGraph
	}
	return nil
}

func NewGraph(startingVertexId string) *Graph {
	v := newVertex(startingVertexId)
	graph := &Graph{
		graph:          make(map[string]*Vertex),
		StartingVertex: v,
		currentVertex:  v,
	}
	graph.graph[startingVertexId] = v
	return graph
}
