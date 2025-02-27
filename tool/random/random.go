package random

import (
	"math/rand"
	"time"
)

const (
	alphaNumCharSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type Random struct {
	rand *rand.Rand
}

func New() *Random {
	return &Random{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (r *Random) GenerateAlphaNum(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = alphaNumCharSet[r.rand.Intn(len(alphaNumCharSet))]
	}
	return string(b)
}
