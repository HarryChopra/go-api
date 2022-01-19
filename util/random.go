package util

import (
	"math/rand"
	"strings"
	"time"
)

var (
	alphabet   = "abcdefghijklmnopqrstuvwxyz"
	currencies = []string{"GBP", "EUR", "USD", "CAD", "AUD"}
)

func init() {
	rand.Seed(time.Now().Unix())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(26)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomName() string {
	return RandomString(8)
}

func RandomCurrency() string {
	n := rand.Intn(len(currencies))
	return currencies[n]
}

func RandomMoney() int64 {
	return RandomInt(1, 1000)
}
