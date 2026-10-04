package cluster

import "testing"

func TestRouter(t *testing.T) {
	nodes := DefaultNodes()

	router := NewRouter(nodes)

	node := router.GetNode("B101")

	if node.ID == "" {
		t.Fatal("router returned an empty node")
	}

	t.Logf("B101 routed to %s (%s)", node.ID, node.Address)
}
