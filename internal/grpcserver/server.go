package grpcserver

import (
	"context"

	"distributed-library-kv/internal/storage"
	pb "distributed-library-kv/proto"
)

type Server struct {
	pb.UnimplementedKVServiceServer

	store *storage.LSMTree
}

func NewServer(store *storage.LSMTree) *Server {
	return &Server{
		store: store,
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

	return &pb.DeleteResponse{
		Success: true,
	}, nil
}
