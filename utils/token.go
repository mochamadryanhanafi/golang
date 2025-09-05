package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateSecureRandomString(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		// Jika crypto/rand gagal, ini adalah masalah serius di level OS.
		// Panic adalah tindakan yang wajar di sini untuk layanan otentikasi.
		panic("could not generate random string: " + err.Error())
	}
	return hex.EncodeToString(bytes)
}
