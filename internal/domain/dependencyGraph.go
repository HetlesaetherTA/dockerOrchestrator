package domain

// A generic asylic directed graph datastructure which lets you keep
// track of dependency chains. Root nodes represent serivces
// and edge nodes are root dependencies.
type DependencyTree[T any] struct {
	Nodes map[string]*Node[T]
}

// Generic tree node
type Node[T any] struct {
	Key   string
	Val   *T
	Edges []Node[T]
}

func (tree *DependencyTree[T]) Get(key string) T {
	return *tree.Nodes[key].Val
}

func (tree *DependencyTree[T]) GetDependencies(key string) []T {
	nodes := tree.Nodes[key]

	if nodes == nil {
		return nil
	}

	tmp := make([]T, 0)

	for _, v := range nodes.Edges {
		tmp = append(tmp, *v.Val)
	}

	return tmp
}

func (tree *DependencyTree[T]) Add() {

}
