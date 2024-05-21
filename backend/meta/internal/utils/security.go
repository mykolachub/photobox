package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func GenerateUUID(filename string) string {
	u := uuid.New()
	data := []byte(u.String() + filename)
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func GetTimeFormatted() string {
	return time.Now().Format("2006-01-02")
}

func GenerateS3FileLocation(fileName, fileExt string) string {
	return fmt.Sprintf("%v/%s%s", GetTimeFormatted(), GenerateUUID(fileName), fileExt)
}
