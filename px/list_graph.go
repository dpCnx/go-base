package main

import (
	"container/list"
	"fmt"
	"reflect"
)

type Vertex struct {
	Value interface{}

	InEdges  []*Edge
	OutEdges []*Edge
}

type Edge struct {
	Weight interface{}

	From *Vertex
	To   *Vertex
}

type ListGraph struct {
	Vertices map[interface{}]*Vertex
	Edges    []*Edge

	Visitor func(v interface{}) bool
}

func NewListGraph(visitor func(v interface{}) bool) *ListGraph {
	return &ListGraph{
		Vertices: make(map[interface{}]*Vertex),
		Edges:    make([]*Edge, 0),
		Visitor:  visitor,
	}
}

func (l *ListGraph) EdgesSize() int {
	return len(l.Vertices)
}

func (l *ListGraph) VerticesSize() int {
	return len(l.Edges)
}

func (l *ListGraph) AddVertex(v interface{}) {
	l.Vertices[v] = &Vertex{
		Value:    v,
		OutEdges: make([]*Edge, 0),
		InEdges:  make([]*Edge, 0),
	}
}

func (l *ListGraph) AddEdge(from, to, weight interface{}) {

	var fromGraph *Vertex
	var toGraph *Vertex
	var ok bool

	fromGraph, ok = l.Vertices[from]
	if !ok {
		fromGraph = &Vertex{
			Value:    from,
			OutEdges: make([]*Edge, 0),
			InEdges:  make([]*Edge, 0),
		}
		l.Vertices[from] = fromGraph
	}

	toGraph, ok = l.Vertices[to]
	if !ok {
		toGraph = &Vertex{
			Value:    to,
			OutEdges: make([]*Edge, 0),
			InEdges:  make([]*Edge, 0),
		}
		l.Vertices[to] = toGraph
	}

	edge := &Edge{
		Weight: weight,
		From:   fromGraph,
		To:     toGraph,
	}

	for _, outEdge := range fromGraph.OutEdges {
		if reflect.DeepEqual(outEdge, edge) {
			return
		}
	}

	fromGraph.OutEdges = append(fromGraph.OutEdges, edge)
	toGraph.InEdges = append(toGraph.InEdges, edge)
	l.Edges = append(l.Edges, edge)
}

func (l *ListGraph) RemoveEdge(from, to interface{}) {

	var fromGraph *Vertex
	var toGraph *Vertex
	var ok bool

	fromGraph, ok = l.Vertices[from]
	if !ok {
		return
	}

	toGraph, ok = l.Vertices[to]
	if !ok {
		return
	}

	edge := &Edge{
		From: fromGraph,
		To:   toGraph,
	}

	var fromGraphIndex int
	for i, outEdge := range fromGraph.OutEdges {
		if reflect.DeepEqual(outEdge.From, edge.From) && reflect.DeepEqual(outEdge.To, edge.To) {
			fromGraphIndex = i
			break
		}
	}

	fromGraph.OutEdges = append(fromGraph.OutEdges[:fromGraphIndex], fromGraph.OutEdges[fromGraphIndex+1:]...)

	var toGraphIndex int
	var isFind bool
	for i, toEdge := range toGraph.InEdges {
		if reflect.DeepEqual(toEdge.From, edge.From) && reflect.DeepEqual(toEdge.To, edge.To) {
			toGraphIndex = i
			isFind = true
			break
		}
	}

	if !isFind {
		return
	}

	toGraph.InEdges = append(toGraph.InEdges[:toGraphIndex], toGraph.InEdges[toGraphIndex+1:]...)

	var edgesIndex int
	for i, e := range l.Edges {
		if reflect.DeepEqual(e.From, edge.From) && reflect.DeepEqual(e.To, edge.To) {
			edgesIndex = i
			break
		}
	}

	l.Edges = append(l.Edges[:edgesIndex], l.Edges[edgesIndex+1:]...)
}

// Bfs 广度优先
func (l *ListGraph) Bfs(begin interface{}) {

	beginVertice, ok := l.Vertices[begin]
	if !ok {
		return
	}

	var visitedVertices []*Vertex
	queue := list.New()

	queue.PushBack(beginVertice)
	visitedVertices = append(visitedVertices, beginVertice)

	for queue.Len() > 0 {
		e := queue.Front()
		v := e.Value.(*Vertex)
		b := l.Visitor(v.Value)
		if b {
			return
		}
		queue.Remove(e)

		for _, edge := range v.OutEdges {

			// 判断定点是否已经存在
			if isContains(visitedVertices, edge.To) {
				return
			}
			queue.PushBack(edge.To)
			visitedVertices = append(visitedVertices, edge.To)
		}
	}

}

// Dfs 深度优先
func (l *ListGraph) Dfs(begin interface{}) {

	beginVertice, ok := l.Vertices[begin]
	if !ok {
		return
	}

	var visitedVertices []*Vertex
	queue := list.New()

	queue.PushBack(beginVertice)
	visitedVertices = append(visitedVertices, beginVertice)

	for queue.Len() > 0 {
		e := queue.Back()
		v := e.Value.(*Vertex)
		queue.Remove(e)

		for _, edge := range v.OutEdges {

			// 判断定点是否已经存在
			if isContains(visitedVertices, edge.To) {
				return
			}
			queue.PushBack(edge.From)
			queue.PushBack(edge.To)
			visitedVertices = append(visitedVertices, edge.To)
			b := l.Visitor(v.Value)
			if b {
				return
			}
		}
	}

}

func isContains(collection []*Vertex, v *Vertex) bool {
	for _, vertex := range collection {
		if reflect.DeepEqual(vertex, v) {
			return true
		}
	}
	return false
}

func main() {
	g := NewListGraph(nil)
	g.AddEdge("A1", "B1", 1)
	g.AddEdge("A1", "B2", 2)
	g.AddEdge("A1", "B3", 3)
	g.AddEdge("B3", "D1", 3)

	/*fmt.Println(g.Vertices)
	fmt.Println(g.Edges)

	g.RemoveEdge("A1", "B2")

	fmt.Println(g.Vertices)
	fmt.Println(g.Edges)

	g.RemoveEdge("A1", "B2")

	fmt.Println(g.Vertices)
	fmt.Println(g.Edges)*/

	g.Visitor = func(v interface{}) bool {
		fmt.Println(v)
		return false
	}

	g.Bfs("A1")

}
