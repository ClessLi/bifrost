package config_graph

type Edge struct {
	sid string
	tid string
}

func newEdge(sid, tid string) Edge {
	return Edge{
		sid: sid,
		tid: tid,
	}
}
