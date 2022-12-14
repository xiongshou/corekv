package utils

import (
	"math/rand"
	"sync"
	"time"
)

var (
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	mu sync.Mutex
)

func Int63n(n int64) int64{
	mu.Lock()
	res := r.Int63n(n)
	mu.Unlock()
	return res
}
func RandN(n int) int {
	mu.Lock()
	res := r.Intn(n)
	mu.Unlock()
	return res
}

func Float64() float64 {
	mu.Lock()
	res := r.Float64()
	mu.Unlock()
	return res
}
