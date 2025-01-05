package utils

import (
	"math/rand"
	"time"
)

func GenerateRandomDate(min, max time.Time) time.Time {
	delta := max.Unix() - min.Unix()
	randTime := min.Add(time.Duration(rand.Int63n(delta)) * time.Second)
	return randTime
}

func GenerateRandomNote() int {
	return rand.Intn(5) + 1
}
