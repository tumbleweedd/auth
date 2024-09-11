package user

import (
	"context"

	userProto "github.com/tumbleweedd/auth_service_proto/gen/go/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) Create(ctx context.Context, req *userProto.CreateUserRequest) (*userProto.CreateUserResponse, error) {
	input := NewCreateUserRequest(req)

	if err := input.Validate(); err != nil {
		return &userProto.CreateUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	createUserOutput, err := s.userUseCase.Create(ctx, &input)
	if err != nil {
		return &userProto.CreateUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	response := NewCreateUserResponse(createUserOutput)

	return response, nil
}

func (s *Server) Update(ctx context.Context, req *userProto.UpdateUserRequest) (*userProto.UpdateUserResponse, error) {
	return &userProto.UpdateUserResponse{}, nil
}

func (s *Server) Delete(ctx context.Context, req *userProto.DeleteUserRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *Server) Get(ctx context.Context, req *userProto.GetUserRequest) (*userProto.GetUserResponse, error) {
	return &userProto.GetUserResponse{}, nil
}

func (s *Server) List(ctx context.Context, req *userProto.ListUsersRequest) (*userProto.ListUsersResponse, error) {
	return &userProto.ListUsersResponse{}, nil
}
