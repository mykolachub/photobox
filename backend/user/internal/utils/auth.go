package utils

import (
	"errors"
	"photobox-user/config"
	"strings"
)

func ValidateBearerHeader(authHeader string) (string, error) {
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return "", errors.New("invalid bearer")
	}

	authType := strings.ToLower(fields[0])
	if authType != config.AuthorizationTypeBearer {
		return "", errors.New("invalid authorization bearer")
	}

	return fields[1], nil
}
