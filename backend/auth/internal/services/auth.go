package services

import (
	"context"
	"encoding/json"
	"errors"
	"photobox-auth/internal/utils"
	"photobox-auth/proto"

	"golang.org/x/oauth2"
)

type AuthService struct {
	GoogleOauthConfig *oauth2.Config
	UserService       proto.UserServiceClient
	Cfg               AuthServiceConfig

	proto.UnimplementedAuthServiceServer
}

type AuthServiceConfig struct {
	JWTSecret string
}

func NewAuthService(googleCfg *oauth2.Config, cfg AuthServiceConfig, userService proto.UserServiceClient) *AuthService {
	return &AuthService{GoogleOauthConfig: googleCfg, Cfg: cfg, UserService: userService}
}

func (s AuthService) GoogleSignup(ctx context.Context, in *proto.GoogleSignupRequest) (*proto.GoogleSignupResponse, error) {
	// TODO: change state
	authStateString := "test"
	url := s.GoogleOauthConfig.AuthCodeURL(authStateString)
	return &proto.GoogleSignupResponse{Url: url}, nil
}

func (s AuthService) GoogleLogin(ctx context.Context, in *proto.GoogleLoginRequest) (*proto.GoogleLoginResponse, error) {
	token, err := s.GoogleOauthConfig.Exchange(context.Background(), in.Code)
	if err != nil {
		return &proto.GoogleLoginResponse{}, err
	}

	client := s.GoogleOauthConfig.Client(context.Background(), token)
	res, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return &proto.GoogleLoginResponse{}, err
	}
	defer res.Body.Close()

	googleUserInfo := &struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"given_name"`
		Picture string `json:"picture"`
	}{}
	if err := json.NewDecoder(res.Body).Decode(googleUserInfo); err != nil {
		return &proto.GoogleLoginResponse{}, err
	}

	// Check if the user exists
	user, err := s.UserService.GetUserByEmail(ctx, &proto.GetUserByEmailRequest{Email: googleUserInfo.Email})
	if err != nil || user.Email == "" {
		user, err = s.UserService.CreateUser(ctx, &proto.CreateUserRequest{
			GoogleId: googleUserInfo.ID,
			Email:    googleUserInfo.Email,
			Username: googleUserInfo.Name,
			Picture:  googleUserInfo.Picture,
		})
		if err != nil {
			return &proto.GoogleLoginResponse{}, errors.New("failed on creating new user")
		}
	}

	// Generate JWT token
	jwtToken, err := utils.GenerateJWTToken(user.Id, user.Email, s.Cfg.JWTSecret)
	if err != nil {
		return &proto.GoogleLoginResponse{}, errors.New("failed on creating jwt token")
	}

	return &proto.GoogleLoginResponse{Token: jwtToken}, nil
}
