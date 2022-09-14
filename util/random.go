package util

import (
	"fmt"
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

func GenerateString(length int) string {
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
	currencies := []string{}
	for k := range supportedCurrencies {
		currencies = append(currencies, k)
	}
	return currencies[rand.Intn(len(currencies))]
}

func RandomEmail() string {
	return fmt.Sprintf("%s@%s.com", GenerateString(5), GenerateString(4))
}
