package grpcserver

import (
	"context"

	"distributed-library-kv/internal/cluster"
	"distributed-library-kv/internal/raft"
	"distributed-library-kv/internal/storage"
	pb "distributed-library-kv/proto"
)

type Server struct {
	pb.UnimplementedKVServiceServer

	store      *storage.LSMTree
	replicator *cluster.Replicator
	raftNode   *raft.Node
}

func NewServer(
	store *storage.LSMTree,
	replicator *cluster.Replicator,
	raftNode *raft.Node,
) *Server {
	return &Server{
		store:      store,
		replicator: replicator,
		raftNode:   raftNode,
	}
}

func (s *Server) Put(
	ctx context.Context,
	req *pb.PutRequest,
) (*pb.PutResponse, error) {

	err := s.store.Put(req.GetKey(), req.GetValue())
	if err != nil {
		return &pb.PutResponse{
			Success: false,
		}, err
	}

	if s.replicator != nil {
		s.replicator.ReplicatePut(
			req.GetKey(),
			req.GetValue(),
		)
	}

	return &pb.PutResponse{
		Success: true,
	}, nil
}

func (s *Server) Get(
	ctx context.Context,
	req *pb.GetRequest,
) (*pb.GetResponse, error) {

	value, err := s.store.Get(req.GetKey())
	if err != nil {
		if err == storage.ErrKeyNotFound {
			return &pb.GetResponse{
				Found: false,
			}, nil
		}

		return nil, err
	}

	return &pb.GetResponse{
		Value: value,
		Found: true,
	}, nil
}

func (s *Server) Delete(
	ctx context.Context,
	req *pb.DeleteRequest,
) (*pb.DeleteResponse, error) {

	err := s.store.Delete(req.GetKey())
	if err != nil {
		return &pb.DeleteResponse{
			Success: false,
		}, err
	}

	if s.replicator != nil {
		s.replicator.ReplicateDelete(req.GetKey())
	}

	return &pb.DeleteResponse{
		Success: true,
	}, nil
}

func (s *Server) ReplicatePut(
	ctx context.Context,
	req *pb.PutRequest,
) (*pb.PutResponse, error) {

	err := s.store.Put(req.GetKey(), req.GetValue())
	if err != nil {
		return &pb.PutResponse{
			Success: false,
		}, err
	}

	return &pb.PutResponse{
		Success: true,
	}, nil
}

func (s *Server) ReplicateDelete(
	ctx context.Context,
	req *pb.DeleteRequest,
) (*pb.DeleteResponse, error) {

	err := s.store.Delete(req.GetKey())
	if err != nil {
		return &pb.DeleteResponse{
			Success: false,
		}, err
	}

	return &pb.DeleteResponse{
		Success: true,
	}, nil
}

// RequestVote handles Raft leader election requests.
func (s *Server) RequestVote(
	ctx context.Context,
	req *pb.VoteRequest,
) (*pb.VoteResponse, error) {

	term, granted := s.raftNode.HandleRequestVote(
		int(req.GetTerm()),
		req.GetCandidateId(),
	)

	return &pb.VoteResponse{
		Term:        int32(term),
		VoteGranted: granted,
	}, nil
}
func (s *Server) Heartbeat(
	ctx context.Context,
	req *pb.HeartbeatRequest,
) (*pb.HeartbeatResponse, error) {

	s.raftNode.BecomeFollower(int(req.GetTerm()))

	return &pb.HeartbeatResponse{
		Success: true,
	}, nil
}
