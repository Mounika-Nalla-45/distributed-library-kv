package raft

import (
	"context"
	"sync"
	"time"

	pb "distributed-library-kv/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type State int

const (
	Follower State = iota
	Candidate
	Leader
)

type Node struct {
	mu sync.Mutex

	ID string

	State State

	CurrentTerm int
	VotedFor    string

	LastHeartbeat time.Time
}

func NewNode(id string) *Node {
	return &Node{
		ID:            id,
		State:         Follower,
		CurrentTerm:   0,
		VotedFor:      "",
		LastHeartbeat: time.Now(),
	}
}

func (n *Node) BecomeCandidate() {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.State = Candidate
	n.CurrentTerm++
	n.VotedFor = n.ID
	n.LastHeartbeat = time.Now()
}

func (n *Node) BecomeLeader() {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.State = Leader
}

func (n *Node) BecomeFollower(term int) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.State = Follower
	n.CurrentTerm = term
	n.VotedFor = ""
	n.LastHeartbeat = time.Now()
}
func (n *Node) HandleRequestVote(term int, candidateID string) (int, bool) {
	n.mu.Lock()
	defer n.mu.Unlock()

	// Reject an old term.
	if term < n.CurrentTerm {
		return n.CurrentTerm, false
	}

	// A newer term makes this node a follower.
	if term > n.CurrentTerm {
		n.CurrentTerm = term
		n.State = Follower
		n.VotedFor = ""
	}

	// Give the vote if we have not voted yet,
	// or if we already voted for this candidate.
	if n.VotedFor == "" || n.VotedFor == candidateID {
		n.VotedFor = candidateID
		n.LastHeartbeat = time.Now()

		return n.CurrentTerm, true
	}

	// Already voted for another candidate.
	return n.CurrentTerm, false
}
func (n *Node) StartElection(nodes []string) bool {
	n.BecomeCandidate()

	n.mu.Lock()
	term := n.CurrentTerm
	candidateID := n.ID
	n.mu.Unlock()

	votes := 1

	for _, address := range nodes {
		if address == "" {
			continue
		}

		ctx, cancel := context.WithTimeout(
			context.Background(),
			2*time.Second,
		)

		conn, err := grpc.NewClient(
			address,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)

		if err != nil {
			cancel()
			continue
		}

		client := pb.NewKVServiceClient(conn)

		response, err := client.RequestVote(
			ctx,
			&pb.VoteRequest{
				Term:        int32(term),
				CandidateId: candidateID,
			},
		)

		conn.Close()
		cancel()

		if err != nil {
			continue
		}

		if response.GetVoteGranted() {
			votes++
		}
	}

	if votes >= 2 {
		n.BecomeLeader()
		return true
	}

	return false
}
func (n *Node) IsLeader() bool {
	n.mu.Lock()
	defer n.mu.Unlock()

	return n.State == Leader
}
func (n *Node) StartHeartbeat(nodes []string) {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			if !n.IsLeader() {
				continue
			}

			for _, address := range nodes {
				go n.sendHeartbeat(address)
			}
		}
	}()
}
func (n *Node) sendHeartbeat(address string) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		2*time.Second,
	)
	defer cancel()

	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return
	}
	defer conn.Close()

	client := pb.NewKVServiceClient(conn)

	n.mu.Lock()
	term := n.CurrentTerm
	n.mu.Unlock()

	_, _ = client.Heartbeat(
		ctx,
		&pb.HeartbeatRequest{
			Term:     int32(term),
			LeaderId: n.ID,
		},
	)
}
