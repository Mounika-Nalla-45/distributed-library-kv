package raft

import "testing"

func TestNewNode(t *testing.T) {
	node := NewNode("node1")

	if node.ID != "node1" {
		t.Errorf("expected node1, got %s", node.ID)
	}

	if node.State != Follower {
		t.Errorf("new node should be follower")
	}

	if node.CurrentTerm != 0 {
		t.Errorf("initial term should be 0")
	}
}

func TestBecomeCandidate(t *testing.T) {
	node := NewNode("node1")

	node.BecomeCandidate()

	if node.State != Candidate {
		t.Errorf("node should become candidate")
	}

	if node.CurrentTerm != 1 {
		t.Errorf("term should become 1")
	}

	if node.VotedFor != "node1" {
		t.Errorf("candidate should vote for itself")
	}
}

func TestBecomeLeader(t *testing.T) {
	node := NewNode("node1")

	node.BecomeCandidate()
	node.BecomeLeader()

	if node.State != Leader {
		t.Errorf("node should become leader")
	}
}

func TestBecomeFollower(t *testing.T) {
	node := NewNode("node1")

	node.BecomeCandidate()
	node.BecomeFollower(5)

	if node.State != Follower {
		t.Errorf("node should become follower")
	}

	if node.CurrentTerm != 5 {
		t.Errorf("expected term 5, got %d", node.CurrentTerm)
	}
}
