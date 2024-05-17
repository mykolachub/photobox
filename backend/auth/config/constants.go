package config

import "time"

var (
	GoogleOauthScopes = []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
		"openid"}
)

var (
	JWTClainsUserId     = "id"
	JWTClainsUserEmail  = "email"
	JWTClainsExpiration = time.Hour * 24
)
