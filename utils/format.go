// (c) Jisin0

package utils

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// Character set to generate a random string
const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var ADMINS []int64

func init() {
	// Create a list of admins from the ADMINS environment variable
	for _, n := range strings.Split(os.Getenv("ADMINS"), " ") {
		num, err := strconv.ParseInt(n, 0, 64)
		if err != nil {
			fmt.Printf("invalid admin id: %v", n)
			continue
		}

		ADMINS = append(ADMINS, num)
	}
}

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
	for _, n := range ADMINS {
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
