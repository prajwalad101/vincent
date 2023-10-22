package utils

import (
	"math/rand"
	"time"
)

func GenerateJobId() int {
	// TODO: make this unique not random
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := r.Intn(10000000)
	return id
}
