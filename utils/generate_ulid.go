package utils

import (
	"github.com/oklog/ulid/v2"
	"math/rand"
	"time"
)

func GenerateULID() string {
	return ulid.MustNew(ulid.Timestamp(time.Now()), rand.New(rand.NewSource(time.Now().UnixNano()))).String()
}
