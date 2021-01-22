package graph

type Edge struct {
	sid     string
	tid     string
	repeats uint
}

func newEdge(sid, tid string) Edge {
	return Edge{
		sid:     sid,
		tid:     tid,
		repeats: 1,
	}
}
