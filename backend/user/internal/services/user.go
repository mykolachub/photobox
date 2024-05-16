package services

import (
	"context"
	"errors"
	"fmt"
	"photobox-user/internal/models/entity"
	"photobox-user/internal/utils"
	"photobox-user/proto"

	"google.golang.org/protobuf/types/known/timestamppb"
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
		return &proto.SignupResponse{}, fmt.Errorf("failed to signup: %s", err)
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
	user, err := s.UserRepo.GetUser(in.UserId)
	if err != nil {
		return &proto.UserResponse{}, errors.New("no such user")
	}

	res := makeUserResponse(user)
	return &res, nil
}

func (s *UserService) GetMe(ctx context.Context, in *proto.GetMeRequest) (*proto.UserResponse, error) {
	return &proto.UserResponse{}, nil
}

func (s *UserService) GetAllUsers(ctx context.Context, in *proto.GetAllUsersRequest) (*proto.GetAllUsersResponse, error) {
	users, err := s.UserRepo.GetAllUsers()
	if err != nil {
		return &proto.GetAllUsersResponse{}, errors.New("failed to get users")
	}

	res := []*proto.UserResponse{}
	for _, v := range users {
		res_user := makeUserResponse(v)
		res = append(res, &res_user)
	}

	return &proto.GetAllUsersResponse{Users: res}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, in *proto.UpdateUserRequest) (*proto.UserResponse, error) {
	updateData := entity.User{
		Email:    in.Email,
		Password: in.Password,
		Username: in.Username,
	}
	if len(in.Password) != 0 {
		hashedPassword, err := utils.HashPassword(in.Password)
		if err != nil {
			return &proto.UserResponse{}, errors.New("failed to set new password")
		}
		updateData.Password = hashedPassword
	}

	user, err := s.UserRepo.UpdateUser(in.UserId, updateData)
	if err != nil {
		return &proto.UserResponse{}, fmt.Errorf("failed to update user %s", err.Error())
	}

	res := makeUserResponse(user)
	return &res, nil
}

func (s *UserService) DeleteUser(ctx context.Context, in *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	user, err := s.UserRepo.DeleteUser(in.UserId)
	if err != nil {
		return &proto.DeleteUserResponse{}, errors.New("failed to delete user")
	}

	res := proto.DeleteUserResponse{UserId: user.ID, Success: true}
	return &res, nil
}

func makeUserResponse(user entity.User) proto.UserResponse {
	return proto.UserResponse{
		Id:          user.ID,
		Email:       user.Email,
		Password:    user.Password,
		Username:    user.Username,
		StorageUsed: user.StorageUsed,
		MaxStorage:  user.MaxStorage,
		CreatedAt:   timestamppb.New(user.CreatedAt),
		UpdatedAt:   timestamppb.New(user.UpdatedAt),
	}
}
