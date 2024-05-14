package services

import (
	"context"
	"photobox-user/proto"
)

type UserService struct {
	UserRepo UserRepo
	proto.UnimplementedUserServiceServer
}

func NewUserService(userRepo UserRepo) *UserService {
	return &UserService{UserRepo: userRepo}
}

func (s *UserService) Signup(ctx context.Context, in *proto.SignupRequest) (*proto.SignupResponse, error) {
	return &proto.SignupResponse{}, nil
}

func (s *UserService) Login(ctx context.Context, in *proto.LoginRequest) (*proto.LoginResponse, error) {
	return &proto.LoginResponse{}, nil
}

func (s *UserService) GetUser(ctx context.Context, in *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	return &proto.GetUserResponse{}, nil
}

func (s *UserService) GetMe(ctx context.Context, in *proto.GetMeRequest) (*proto.UserResponse, error) {
	return &proto.UserResponse{}, nil
}

func (s *UserService) GetAllUsers(ctx context.Context, in *proto.GetAllUsersRequest) (*proto.GetAllUsersResponse, error) {
	return &proto.GetAllUsersResponse{}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, in *proto.UpdateUserRequest) (*proto.GetUserResponse, error) {
	return &proto.GetUserResponse{}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, in *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	return &proto.DeleteUserResponse{}, nil
}
