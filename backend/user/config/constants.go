package config

var (
	AuthorizationHeader     = "authorization"
	AuthorizationTypeBearer = "bearer"

	PayloadUserId    = "payload_user_id"
	PayloadUserEmail = "payload_user_email"

	JWTClainsUserId     = "id"
	JWTClainsUserEmail  = "email"
	JWTClainsExpiration = "exp"
)
