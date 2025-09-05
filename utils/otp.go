package utils

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
)

// GenerateOTP menghasilkan OTP numerik acak dengan panjang yang ditentukan.
func GenerateOTP(length int) string {
	b := make([]byte, length)
	n, err := io.ReadAtLeast(rand.Reader, b, length)
	if n != length || err != nil {
		// Fallback ke metode yang lebih sederhana jika crypto/rand gagal
		return fallbackOTP(length)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

// fallbackOTP adalah generator cadangan jika crypto/rand tidak tersedia.
func fallbackOTP(length int) string {
	max := new(big.Int)
	max.SetString(string(make([]byte, length+1)), 10) // 10^length
	max.Bytes()[0] = '1'

	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		panic(err) // Seharusnya tidak pernah terjadi
	}

	return fmt.Sprintf("%0*d", length, n)
}
