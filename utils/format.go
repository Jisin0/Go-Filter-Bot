// (c) Jisin0

package utils

import (
	"math/rand"
	"time"

	"github.com/Jisin0/Go-Filter-Bot/utils/config"
)

// Character set to generate a random string
const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandString(length int) string {
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}

// IsAdmin returns wether the user is a global admin(from env).
func IsAdmin(u int64) bool {
	for _, n := range config.Admins {
		if u == n {
			return true
		}
	}

	return false
}

// Checks if a string slice contains an item.
func Contains(l []string, v string) bool {
	for _, i := range l {
		if i == v {
			return true
		}
	}

	return false
}
