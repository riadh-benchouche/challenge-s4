package utils

import (
	"math/rand"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
)

func GenerateULID() string {
	return ulid.MustNew(ulid.Timestamp(time.Now()), rand.New(rand.NewSource(time.Now().UnixNano()))).String()
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateAssociationCode() string {
	length := rand.Intn(16) + 5 // Longueur aléatoire entre 5 et 20

	var builder strings.Builder
	for i := 0; i < length; i++ {
		builder.WriteByte(charset[rand.Intn(len(charset))])
	}

	return builder.String()
}
