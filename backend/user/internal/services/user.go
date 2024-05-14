package services

import (
	"context"
	"errors"
	"fmt"
	"photobox-user/internal/models/entity"
	"photobox-user/internal/utils"
	"photobox-user/proto"
)

type UserService struct {
	UserRepo UserRepo
	proto.UnimplementedUserServiceServer

	cfg UserServiceConfig
}

type UserServiceConfig struct {
	JwtSecret string
}

func NewUserService(r UserRepo, c UserServiceConfig) *UserService {
	return &UserService{UserRepo: r, cfg: c}
}

func (s *UserService) Signup(ctx context.Context, in *proto.SignupRequest) (*proto.SignupResponse, error) {
	user, _ := s.UserRepo.GetUserByEmail(in.Email)

	if user.Email != "" {
		return &proto.SignupResponse{}, errors.New("user already exists")
	}

	hash, err := utils.HashPassword(in.Password)
	if err != nil {
		return &proto.SignupResponse{}, errors.New("failed to hash password")
	}

	_, err = s.UserRepo.CreateUser(entity.User{
		Email:    in.Email,
		Username: in.Username,
		Password: hash,
	})
	if err != nil {
		return &proto.SignupResponse{}, fmt.Errorf("failed to sigup: %s", err)
	}

	res := proto.SignupResponse{Success: true}
	return &res, nil
}

func (s *UserService) Login(ctx context.Context, in *proto.LoginRequest) (*proto.LoginResponse, error) {
	// get user
	user, err := s.UserRepo.GetUserByEmail(in.Email)
	if err != nil {
		return &proto.LoginResponse{}, errors.New("invalid credentials")
	}
	fmt.Printf("user: %+v\n", user)

	// compare password and hash
	match, err := utils.ComparePasswordAndHash(in.Password, user.Password)
	if err != nil || !match {
		return &proto.LoginResponse{}, errors.New("invalid credentials")
	}
	fmt.Printf("match: %v\n", match)

	// generate token
	token, err := utils.GenerateJWTToken(user.ID, user.Email, s.cfg.JwtSecret)
	if err != nil {
		return &proto.LoginResponse{}, errors.New("failed to generate token")
	}

	res := proto.LoginResponse{Token: token}
	return &res, nil
}

func (s *UserService) GetUser(ctx context.Context, in *proto.GetUserRequest) (*proto.UserResponse, error) {
	return &proto.UserResponse{}, nil
}

func (s *UserService) GetMe(ctx context.Context, in *proto.GetMeRequest) (*proto.UserResponse, error) {
	return &proto.UserResponse{}, nil
}

func (s *UserService) GetAllUsers(ctx context.Context, in *proto.GetAllUsersRequest) (*proto.GetAllUsersResponse, error) {
	return &proto.GetAllUsersResponse{}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, in *proto.UpdateUserRequest) (*proto.UserResponse, error) {
	return &proto.UserResponse{}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, in *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	return &proto.DeleteUserResponse{}, nil
}
