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
func TestRequestVoteRejectsSecondCandidate(t *testing.T) {
	node := NewNode("node1")

	term, granted := node.HandleRequestVote(1, "node2")

	if !granted {
		t.Fatal("first vote should be granted")
	}

	if term != 1 {
		t.Fatalf("expected term 1, got %d", term)
	}

	_, granted = node.HandleRequestVote(1, "node3")

	if granted {
		t.Fatal("node should not vote for two candidates in the same term")
	}
}

func TestRequestVoteHigherTerm(t *testing.T) {
	node := NewNode("node1")

	node.HandleRequestVote(1, "node2")

	term, granted := node.HandleRequestVote(2, "node3")

	if !granted {
		t.Fatal("vote should be granted for higher term")
	}

	if term != 2 {
		t.Fatalf("expected term 2, got %d", term)
	}
}

func TestAppendEntriesMakesFollower(t *testing.T) {
	node := NewNode("node1")

	node.BecomeCandidate()

	term, success := node.HandleAppendEntries(
		2,
		"node2",
		nil,
		0,
	)

	if !success {
		t.Fatal("AppendEntries should succeed")
	}

	if term != 2 {
		t.Fatalf("expected term 2, got %d", term)
	}

	if node.IsLeader() {
		t.Fatal("node should not remain leader")
	}
}
