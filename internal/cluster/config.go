package cluster

type Node struct {
	ID      string
	Address string
}

type Config struct {
	Self  Node
	Nodes []Node
}

func DefaultNodes() []Node {
	return []Node{
		{
			ID:      "node1",
			Address: "localhost:50051",
		},
		{
			ID:      "node2",
			Address: "localhost:50052",
		},
		{
			ID:      "node3",
			Address: "localhost:50053",
		},
	}
}
