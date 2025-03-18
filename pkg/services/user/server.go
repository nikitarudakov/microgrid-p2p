package user

import (
	"context"
	"fmt"
	"github.com/nikitarudakov/microenergy/internal/pb"
	"google.golang.org/protobuf/types/known/emptypb"
	"sync"
)

type Server struct {
	users []*pb.User
	mu    sync.Mutex
	pb.UserManagementServer
}

func (s *Server) FetchAllUsers(_ context.Context, _ *emptypb.Empty) (*pb.UserList, error) {
	return &pb.UserList{Users: s.users}, nil
}

func (s *Server) FetchUser(_ context.Context, in *pb.FetchUserInput) (*pb.User, error) {
	if in.Id != nil {
		return s.fetchUserByID(in.GetId())
	}

	return nil, nil
}

func (s *Server) RegisterUser(_ context.Context, in *pb.RegisterUserInput) (*pb.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user := &pb.User{
		Id:        in.Id,
		FirstName: in.FirstName,
		LastName:  in.LastName,
	}

	s.users = append(s.users, user)

	return user, nil
}

func (s *Server) fetchUserByID(id string) (*pb.User, error) {
	for _, u := range s.users {
		if u.Id == id {
			return u, nil
		}
	}

	return nil, fmt.Errorf("user with id [%s] was not found", id)
}
