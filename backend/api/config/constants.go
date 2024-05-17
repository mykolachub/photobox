package config

import "time"

var (
	JWTClainsUserId     = "id"
	JWTClainsUserEmail  = "email"
	JWTClainsExpiration = time.Hour * 24 // One day
)

var (
	AuthorizationHeader     = "authorization"
	AuthorizationTypeBearer = "bearer"
)

var (
	PayloadUserId    = "payload_user_id"
	PayloadUserEmail = "payload_user_email"
)
