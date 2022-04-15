package graph

import "strings"

type Vertex struct {
	id     string
	edges  []Edge
	degree uint
}

func (v Vertex) Size() uint {
	return v.degree
}

func (v Vertex) GetEdgeByIndex(idx int) Edge {
	if idx > int(v.degree) {
		return Edge{}
	}

	return v.edges[idx]
}

func (v *Vertex) AddEdge(id string) error {
	if strings.EqualFold(v.id, id) {
		return ErrVertexRepeat
	}
	for i := 0; i < int(v.degree); i++ {
		if strings.EqualFold(v.edges[i].tid, id) {
			// return ErrEdgeRepeat
			v.edges[i].repeats++

			return nil
		}
	}
	v.edges = append(v.edges, newEdge(v.id, id))
	v.degree++

	return nil
}

func (v *Vertex) DelEdge(id string) {
	if id == v.id {
		return
	}
	for i := 0; i < int(v.degree); i++ {
		if strings.EqualFold(v.edges[i].tid, id) {
			if v.edges[i].repeats > 1 {
				v.edges[i].repeats--

				return
			}
			v.edges = append(v.edges[:i], v.edges[i+1:]...)
			v.degree--

			return
		}
	}
}

func newVertex(id string) *Vertex {
	return &Vertex{
		id:     id,
		edges:  make([]Edge, 0),
		degree: 0,
	}
}
