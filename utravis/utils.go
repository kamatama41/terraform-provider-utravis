package utravis

import (
	"crypto/sha256"
	"encoding/base64"
)

func hashString(str string) string {
	hash := sha256.Sum256([]byte(str))
	return base64.StdEncoding.EncodeToString(hash[:])
}
