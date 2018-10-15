package tbstring

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandomString(wordLength int) string {
	b := make([]rune, wordLength)
	max := len(letterRunes)
	for i := range b {
		randomIndex := rand.Intn(max)
		b[i] = letterRunes[randomIndex]
	}
	return string(b)
}
