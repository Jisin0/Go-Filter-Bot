// (c) Jisin0

package utils

import (
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var Dst []byte = make([]byte, 25)

// Character set to generate a random string
const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandString(length int) string {
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func GetAdmins() []int64 {
	// Create a list of admins from the ADMINS environment variable
	var nums []int64

	for _, n := range strings.Split(os.Getenv("ADMINS"), " ") {
		num, _ := strconv.ParseInt(n, 0, 64)
		nums = append(nums, num)
	}

	return nums
}
