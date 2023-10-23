package util

import (
	"math/rand"
	"strconv"
	"time"
)

func GenerateJobId() string {
	// TODO: make this unique not random
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := r.Intn(10000000)
	return strconv.Itoa(id)
}
