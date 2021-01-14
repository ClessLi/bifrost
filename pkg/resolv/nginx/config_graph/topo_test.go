package config_graph

import "testing"

func TestTopo_CircleVertexIndex(t *testing.T) {
	graph := NewGraph("nginx.conf")
	graph.AddEdge("nginx.conf", "./conf.d/test1.conf")
	graph.AddEdge("nginx.conf", "./conf.d/test2.conf")
	graph.AddEdge("nginx.conf", "./conf.d/test3.conf")
	graph.AddEdge("nginx.conf", "./conf.d/test4.conf")
	graph.SetCurrentVertex("./conf.d/test1.conf")
	graph.AddEdgeFromCurrentVertex("./conf.d/test/1.conf")
	graph.SetCurrentVertex("./conf.d/test2.conf")
	graph.AddEdgeFromCurrentVertex("./conf.d/location.conf")
	graph.SetCurrentVertex("./conf.d/test3.conf")
	graph.AddEdgeFromCurrentVertex("./conf.d/location.conf")
	graph.SetCurrentVertex("./conf.d/test4.conf")
	graph.AddEdgeFromCurrentVertex("./conf.d/location.conf")
	//graph.AddEdge("./conf.d/test1.conf", "./conf.d/test/1.conf")
	//graph.AddEdge("./conf.d/test1.conf", "./conf.d/location.conf")
	//graph.AddEdge("./conf.d/test2.conf", "./conf.d/location.conf")
	//graph.AddEdge("./conf.d/test3.conf", "./conf.d/location.conf")
	//graph.AddEdge("./conf.d/test4.conf", "./conf.d/location.conf")
	//graph.AddEdge("./conf.d/test/1.conf", "nginx.test.conf")
	//graph.AddEdge("./conf.d/test/1.conf", "nginx.conf")
	graph.AddEdge("nginx.test.conf", "./conf.d/test1.conf")
	err := graph.TopoLogicalSortByKahn()
	if err != nil {
		t.Fatal(err)
	}
}
