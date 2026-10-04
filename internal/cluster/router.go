package cluster

import "hash/fnv"

type Router struct {
	nodes []Node
}

func NewRouter(nodes []Node) *Router {
	return &Router{
		nodes: nodes,
	}
}

func (r *Router) GetNode(key string) Node {
	if len(r.nodes) == 0 {
		return Node{}
	}

	hash := fnv.New32a()
	_, _ = hash.Write([]byte(key))

	index := int(hash.Sum32()) % len(r.nodes)

	return r.nodes[index]
}
