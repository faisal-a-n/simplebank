package util

import (
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenerateRandomInt(max, min int64) int64 {
	return min + rand.Int63n(max-min+1)
}

const alphabets string = "adcdefghijklmnopqrstuvwzyz"

func GenerateName(length int) string {
	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteByte(alphabets[rand.Intn(len(alphabets))])
	}
	return sb.String()
}

func GenerateAmount() int64 {
	return GenerateRandomInt(1000, 100)
}

func GenerateCurrency() string {
	currencies := []string{"INR", "USD", "EUR", "CAD", "YEN"}
	return currencies[rand.Intn(len(currencies))]
}
