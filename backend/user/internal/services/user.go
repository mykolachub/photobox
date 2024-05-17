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

func (s *UserService) CreateUser(ctx context.Context, in *proto.CreateUserRequest) (*proto.UserResponse, error) {
	user, _ := s.UserRepo.GetUserByEmail(in.Email)

	if user.Email != "" {
		return &proto.UserResponse{}, errors.New("user already exists")
	}

	hash, err := utils.HashPassword(in.Password)
	if err != nil {
		return &proto.UserResponse{}, errors.New("failed to hash password")
	}

	user, err = s.UserRepo.CreateUser(entity.User{
		GoogleID: in.GoogleId,
		Email:    in.Email,
		Password: hash,
		Username: in.Username,
		Picture:  in.Picture,
	})
	if err != nil {
		return &proto.UserResponse{}, fmt.Errorf("failed to signup: %s", err)
	}

	res := makeUserResponse(user)
	return &res, nil
}

func (s *UserService) GetUser(ctx context.Context, in *proto.GetUserRequest) (*proto.UserResponse, error) {
	user, err := s.UserRepo.GetUser(in.Id)
	if err != nil {
		return &proto.UserResponse{}, errors.New("no such user")
	}

	res := makeUserResponse(user)
	return &res, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, in *proto.GetUserByEmailRequest) (*proto.UserResponse, error) {
	user, err := s.UserRepo.GetUserByEmail(in.Email)
	if err != nil {
		return &proto.UserResponse{}, errors.New("no such user")
	}

	res := makeUserResponse(user)
	return &res, nil
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
		Picture:  in.Picture,
	}
	if len(in.Password) != 0 {
		hashedPassword, err := utils.HashPassword(in.Password)
		if err != nil {
			return &proto.UserResponse{}, errors.New("failed to set new password")
		}
		updateData.Password = hashedPassword
	}

	user, err := s.UserRepo.UpdateUser(in.Id, updateData)
	if err != nil {
		return &proto.UserResponse{}, fmt.Errorf("failed to update user %s", err.Error())
	}

	res := makeUserResponse(user)
	return &res, nil
}

func (s *UserService) DeleteUser(ctx context.Context, in *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	user, err := s.UserRepo.DeleteUser(in.Id)
	if err != nil {
		return &proto.DeleteUserResponse{}, errors.New("failed to delete user")
	}

	res := proto.DeleteUserResponse{Id: user.ID, Success: true}
	return &res, nil
}

func makeUserResponse(user entity.User) proto.UserResponse {
	return proto.UserResponse{
		Id:          user.ID,
		GoogleId:    user.GoogleID,
		Email:       user.Email,
		Password:    user.Password,
		Username:    user.Username,
		Picture:     user.Picture,
		StorageUsed: user.StorageUsed,
		MaxStorage:  user.MaxStorage,
		CreatedAt:   timestamppb.New(user.CreatedAt),
		UpdatedAt:   timestamppb.New(user.UpdatedAt),
	}
}
